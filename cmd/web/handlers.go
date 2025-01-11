package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
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

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new client has connected")

	// Set CORS headers before sending anything to client
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// TODO: Configure properly in production
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	rc := http.NewResponseController(w)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	clientGone := r.Context().Done()

	for {
		select {
		case <-clientGone:
			fmt.Println("client has disconnected")
			return
		case <-ticker.C:
			data := fmt.Sprintf(`{"message": "%s"}`, time.Now().Format("15:04:05"))
			fmt.Fprintf(w, "event:ticker\ndata:%s\n\n", data)
			rc.Flush()
		}
	}
}
