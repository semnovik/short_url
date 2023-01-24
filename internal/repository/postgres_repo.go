package repository

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"short_url/configs"
)

type PostgresRepo struct {
	Conn *sql.DB
}

var ErrAlreadyExists = fmt.Errorf("already exists")

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{Conn: db}
}

func (r *PostgresRepo) Add(url string) (string, error) {
	for {
		uuid := GenUUID()

		_, err := r.Conn.Exec(`
			INSERT INTO urls (original_url, short_url)
			VALUES ($1, $2)
		`, url, configs.Config.BaseURL+"/"+uuid)
		if err != nil {
			return "", err
		}

		if errors.Is(err, ErrAlreadyExists) {
			continue
		}

		return uuid, nil
	}
}

func (r *PostgresRepo) Get(uuid string) (string, error) {
	var originalUrl string
	err := r.Conn.QueryRow("SELECT original_url FROM urls WHERE short_url=$1", configs.Config.BaseURL+"/"+uuid).Scan(&originalUrl)
	if err != nil {
		return "", err
	}

	return originalUrl, nil
}

func (r *PostgresRepo) AddByUser(userID, originalURL, shortURL string) {
	_, err := r.Conn.Exec(`
		UPDATE urls 
		SET user_uuid=$1
		WHERE short_url=$2 AND original_url=$3
	`, userID, shortURL, originalURL)
	if err != nil {
		log.Print(err)
	}
}

func (r *PostgresRepo) AllUsersURLS(userID string) []URLObj {
	rows, err := r.Conn.Query(`
		SELECT short_url, original_url
		FROM urls
		WHERE user_uuid=$1
	`, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var (
		urlsByUser []URLObj
		urlByUser  URLObj
	)

	for rows.Next() {
		err := rows.Scan(&urlByUser.ShortURL, &urlByUser.OriginalURL)
		if err != nil {
			return nil
		}
		urlsByUser = append(urlsByUser, urlByUser)
	}
	return urlsByUser
}

func (r *PostgresRepo) IsUserExist(userID string) bool {
	var uuidFromDb string
	row := r.Conn.QueryRow(`
		SELECT user_uuid 
		FROM urls 
		WHERE user_uuid=$1
		LIMIT 1
		`, userID)

	err := row.Scan(&uuidFromDb)
	if err != nil {
		return false
	}
	if uuidFromDb == userID {
		return true
	}
	return false
}

func (r *PostgresRepo) Ping() error {
	if r.Conn == nil {
		return errors.New("something wrong with DB-connection")
	}
	return r.Conn.Ping()
}
