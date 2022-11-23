package main

import (
	"log"
	"net/http"
	"short_url/internal/app/handlers"
	"short_url/internal/app/server"
	"short_url/internal/app/storage"
)

func main() {

	var URLStorage []string

	Repository := storage.NewRepository(URLStorage)
	Server := server.NewServer(Repository)
	Handler := handlers.NewHandler(Server)

	err := http.ListenAndServe(":8080", Handler.InitRouter())
	if err != nil {
		log.Fatal(err)
	}
}
