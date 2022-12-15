package main

import (
	"log"
	"os"
	"short_url/internal/repository"
	"short_url/internal/server"
)

func main() {

	repo := repository.NewURLRepository()
	srv := server.NewShorterSrv(repo)

	srv.Addr = os.Getenv("SERVER_ADDRESS")

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
