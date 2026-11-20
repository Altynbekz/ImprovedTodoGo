package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api", http.StripPrefix("/api", http.FileServer(http.Dir("."))))

	server := http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	log.Println("The base if the server is created")
	log.Fatal(server.ListenAndServe())
}
