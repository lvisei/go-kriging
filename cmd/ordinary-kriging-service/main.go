package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/train", trainHandler)
	http.HandleFunc("/grid", gridHandler)
	server := &http.Server{
		Addr:           ":8888",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", "ordinary-kriging")
}

func trainHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:
}

func gridHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:
}
