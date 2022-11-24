package main

import (
	"log"
	"net/http"
	"short_url/internal/app/handlers"
	"short_url/internal/app/repository"
	"short_url/internal/app/service"
)

func main() {

	var URLStorage []string

	Repository := repository.NewRepository(URLStorage)
	Server := service.NewServer(Repository)
	Handler := handlers.NewHandler(Server)

	err := http.ListenAndServe(":8080", Handler.InitRouter())
	if err != nil {
		log.Fatal(err)
	}
}
