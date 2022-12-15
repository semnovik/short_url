package main

import (
	"log"
	"short_url/configs"
	"short_url/internal/repository"
	"short_url/internal/server"
)

func main() {

	repo := repository.NewURLRepository()
	srv := server.NewShorterSrv(repo)

	srv.Addr = configs.Config.ServerAddress

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
