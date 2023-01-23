package main

import (
	"log"
	"short_url/configs"
	"short_url/internal/repository"
	"short_url/internal/server"
)

func main() {
	// iter11

	configs.InitFlags()

	dbConn, err := repository.NewPostgresRepo()
	if dbConn != nil {
		defer dbConn.Conn.Close()
	}
	if err != nil {
		panic(err)
	}

	repo := repository.NewRepo(dbConn)
	srv := server.NewShorterSrv(repo)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
