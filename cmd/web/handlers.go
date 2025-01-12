package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/fsnotify/fsnotify"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	templateData := struct {
		Version   string
		BuildTime string
		Host      string
		Port      int
	}{
		Version:   version,
		BuildTime: buildTime,
		Host:      host,
		Port:      port,
	}

	err = ts.Execute(w, templateData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("new client has connected")

	// Set CORS headers before sending anything to client
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// TODO: Configure properly in production
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	rc := http.NewResponseController(w)

	clientGone := r.Context().Done()

	// TODO: refactor - move file watching logic to goroutine
	file, err := os.Open(watchedFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//nolint:all
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Set initial read position to the end of file
	readPosition := stat.Size()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err.Error())
		return
	}
	//nolint:all
	defer watcher.Close()

	err = watcher.Add(watchedFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for {
		select {
		case <-clientGone:
			log.Println("client has disconnected")
			return
		case event := <-watcher.Events:
			log.Println(event)

			stat, err := file.Stat()
			if err != nil {
				log.Println(err.Error())
				return
			}

			// Only read from file if it's size has increased
			size := stat.Size()
			if size > readPosition {
				// Make buffer of size just enough to read all new content
				buf := make([]byte, size-readPosition)
				_, err = file.ReadAt(buf, readPosition)
				if err != nil && err.Error() != "EOF" {
					log.Printf("error reading from log file: %s", err.Error())
					return
				}

				_, err = fmt.Fprintf(w, "event:log_updated\ndata:%s\n\n", buf)
				if err != nil {
					log.Println(err.Error())
				}

				// Send data to client
				err = rc.Flush()
				if err != nil {
					log.Println(err.Error())
				}
			}

			// Again, set initial read position to the end of file
			readPosition = size
		case errors := <-watcher.Errors:
			log.Println(errors)
		}
	}
}
