package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"short_url/configs"
)

func NewPostgresDB() *pgx.Conn {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, configs.Config.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
