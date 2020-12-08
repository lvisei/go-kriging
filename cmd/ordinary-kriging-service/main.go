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
	http.HandleFunc("/grid-png", gridPngHandler)
	server := &http.Server{
		Addr:           ":8888",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", "ordinary-kriging")
}

func trainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", "ordinary-kriging")
	// TODO:
}

func gridHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", "ordinary-kriging")
	// TODO:
}

func gridPngHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", "ordinary-kriging")
	// TODO:
}
