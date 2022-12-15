package main

import (
	"log"
	"short_url/internal/repository"
	"short_url/internal/server"
)

func main() {

	repo := repository.NewURLRepository()
	srv := server.NewShorterSrv(repo)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
