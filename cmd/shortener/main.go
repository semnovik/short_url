package main

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"log"
	"short_url/configs"
	"short_url/internal/repository"
	"short_url/internal/server"
)

func main() {
	// iter14.
	configs.InitFlags()

	var db *sql.DB
	var err error

	if configs.Config.DatabaseDSN != "" {
		db, err = sql.Open("pgx", configs.Config.DatabaseDSN)
		if err != nil {
			log.Panicln(err)
		}
		defer db.Close()

		err = goose.Up(db, configs.Config.MigrationsDir, goose.WithAllowMissing())
		if err != nil {
			log.Print("Cannot make migrations!", err)
		}
	}

	repo := repository.NewRepo(db)
	srv := server.NewShorterSrv(repo)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
