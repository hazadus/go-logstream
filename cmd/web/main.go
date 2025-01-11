package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	port = 8000
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
	log.Printf("Waiting for connections on port %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatalf("failed to start server: %s", err.Error())
	}
}
