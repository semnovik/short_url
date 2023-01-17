package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"short_url/configs"
	"short_url/internal/repository"
	"short_url/internal/server"
)

func main() {
	configs.InitFlags()
	ctx := context.Background()

	var dbConn *pgx.Conn
	var err error

	if _, exist := os.LookupEnv("DATABASE_DSN"); exist {
		dbConn, err = pgx.Connect(ctx, configs.Config.DatabaseDSN)
		if err != nil {
			panic(err)
		}
		defer dbConn.Close(ctx)

		err = dbConn.Ping(ctx)
		if err != nil {
			panic(err)
		}
	}

	repo := repository.NewRepo(dbConn)
	srv := server.NewShorterSrv(repo)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
