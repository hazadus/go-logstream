package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
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

	newData := make(chan []byte)
	defer close(newData)

	done := make(chan struct{})
	defer close(done)

	go watchFile(watchedFilePath, newData, done)

	for {
		select {
		case <-clientGone:
			log.Println("client has disconnected")
			// Tell goroutine to stop:
			done <- struct{}{}
			return

		case buf := <-newData:
			_, err := fmt.Fprintf(w, "event:log_updated\ndata:%s\n\n", buf)
			if err != nil {
				log.Println(err.Error())
			}

			// Send data to client
			err = rc.Flush()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}
