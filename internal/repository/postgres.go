package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"short_url/configs"
)

func NewPostgresDB() *pgx.Conn {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, configs.Config.DatabaseDSN)
	if err != nil {
		fmt.Print(err)
	}

	return db
}
