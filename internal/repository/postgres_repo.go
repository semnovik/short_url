package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
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
		`, url, uuid)
		if err != nil {
			return "", err
		}

		if errors.Is(err, ErrAlreadyExists) {
			continue
		}

		return uuid, nil
	}
}

func (r *PostgresRepo) Get(uuid string) (string, bool, error) {
	var originalURL string
	var isDeleted bool

	err := r.Conn.QueryRow("SELECT original_url, is_deleted FROM urls WHERE short_url=$1", uuid).Scan(&originalURL, &isDeleted)
	if err != nil {
		return "", false, err
	}

	return originalURL, isDeleted, nil
}

func ErrAlreadyExist(err error) bool {
	newErr, ok := err.(*pgconn.PgError)
	return ok && newErr.Code == pgerrcode.UniqueViolation
}

func (r *PostgresRepo) AddByUser(userID, originalURL string) (string, error) {

	uuid := GenUUID()

	params := []interface{}{
		originalURL,
		uuid,
		userID,
	}

	query := `INSERT INTO 
				  urls (original_url, short_url, user_uuid)
		          VALUES ($1, $2, $3)`

	_, err1 := r.Conn.Exec(query, params...)
	if ErrAlreadyExist(err1) {
		uuidFromRepo, err2 := r.GetShortByOriginal(originalURL)
		if err2 != nil {
			return "", err2
		}
		return uuidFromRepo, err1
	}

	return uuid, nil
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
		urlByUser.ShortURL = configs.Config.BaseURL + "/" + urlByUser.ShortURL
		urlsByUser = append(urlsByUser, urlByUser)
	}

	err = rows.Err()
	if err != nil {
		return nil
	}

	return urlsByUser
}

func (r *PostgresRepo) IsUserExist(userID string) bool {
	var uuidFromDB string

	row := r.Conn.QueryRow(`
		SELECT user_uuid 
		FROM urls 
		WHERE user_uuid=$1
		LIMIT 1
		`, userID)

	err := row.Scan(&uuidFromDB)
	if err != nil {
		return false
	}
	if uuidFromDB == userID {
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

func (r *PostgresRepo) GetShortByOriginal(originalURL string) (string, error) {
	var uuid string
	query := `SELECT short_url
			  FROM urls
			  WHERE original_url=$1`

	row := r.Conn.QueryRow(query, originalURL)
	err := row.Scan(&uuid)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (r *PostgresRepo) DeleteByUUID(uuid []string, userID string) {
	query := `UPDATE urls SET is_deleted=TRUE WHERE short_url=ANY($1) and user_uuid=$2`

	_, _ = r.Conn.Exec(query, uuid, userID)
}
