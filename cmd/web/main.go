package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	defaultHost          = "localhost"
	defaultPort          = 8000
	defaultWatchFilePath = "./genlog.log"
)

type config struct {
	host          string
	port          int
	watchFilePath string
}

type application struct {
	cfg *config
}

// Устанавливается при сборке при помощи параметра
// go build -ldflags='-X main.buildTime=${current_time}' ...
var buildTime string
var version string

func main() {
	hostFlag := flag.String("host", defaultHost, "Host name")
	portFlag := flag.Int("port", defaultPort, "Port to start app on")
	pathFlag := flag.String("path", defaultWatchFilePath, "Log file to watch, full path")
	flag.Parse()

	if !fileExists(*pathFlag) {
		log.Printf("file not found: %s", *pathFlag)
		os.Exit(1)
	}

	app := &application{
		cfg: &config{
			host:          *hostFlag,
			port:          *portFlag,
			watchFilePath: *pathFlag,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.indexHandler)
	mux.HandleFunc("/events", app.sseHandler)

	log.Printf("Logstream %s, built on %s", version, buildTime)
	log.Printf("Watching file: %s", app.cfg.watchFilePath)
	log.Printf("Waiting for connections on %s:%d\n", app.cfg.host, app.cfg.port)

	//nolint:all
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", app.cfg.host, app.cfg.port), mux)
	if err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}
