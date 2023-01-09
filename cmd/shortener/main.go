package main

import (
	"log"
	"short_url/configs"
	"short_url/internal/repository"
	"short_url/internal/server"
)

func main() {
	configs.InitFlags()
	repo := repository.NewRepo()

	srv := server.NewShorterSrv(repo)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
