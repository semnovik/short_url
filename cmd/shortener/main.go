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
	// iter11

	configs.InitFlags()
	ctx := context.Background()

	var dbConn *pgx.Conn
	var err error

	if configs.Config.DatabaseDSN != "" {
		dbConn, err = pgx.Connect(ctx, configs.Config.DatabaseDSN)

		switch {
		case err != nil:
			log.Print("DB not configured")
		default:
			defer dbConn.Close(ctx)

			err = dbConn.Ping(ctx)
			if err != nil {
				panic(err)
			}
			log.Print("DB successfully configured")
		}
	} else {
		log.Print("DB not configured")
	}

	repo := repository.NewRepo(dbConn)
	srv := server.NewShorterSrv(repo)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
