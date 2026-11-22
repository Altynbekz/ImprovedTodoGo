package main

import (
	"Altyn/projectOne/internal"
	"Altyn/projectOne/internal/api"
	"log"
	"net/http"
)

func main() {

	cfg := &api.ApiConfig{}

	server := http.Server{
		Addr:    ":8082",
		Handler: internal.Middleware(internal.NewRouter(cfg)),
	}

	log.Println("The base if the server is created")
	log.Fatal(server.ListenAndServe())
}
