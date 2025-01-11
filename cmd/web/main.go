package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	host            = "localhost"
	port            = 8000
	watchedFilePath = "./genlog.log"
)

// Устанавливается при сборке при помощи параметра
// go build -ldflags='-X main.buildTime=${current_time}' ...
var buildTime string
var version string

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/events", sseHandler)

	log.Printf("Logstream %s, built on %s", version, buildTime)
	log.Printf("Watching file: %s", watchedFilePath)
	log.Printf("Waiting for connections on port %d\n", port)

	//nolint:all
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}
