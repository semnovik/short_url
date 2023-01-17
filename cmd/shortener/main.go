package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"short_url/configs"
	"short_url/internal/repository"
	"short_url/internal/server"
)

func main() {
	configs.InitFlags()
	ctx := context.Background()

	repo := repository.NewRepo(nil)

	if configs.Config.DatabaseDSN != "" {
		conn, err := pgx.Connect(ctx, configs.Config.DatabaseDSN)
		if err != nil {
			panic(err)
		}
		defer conn.Close(ctx)

		err = conn.Ping(ctx)
		if err != nil {
			panic(err)
		}

		repo = repository.NewRepo(conn)
	}

	srv := server.NewShorterSrv(repo)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
