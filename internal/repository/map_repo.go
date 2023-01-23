package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type URLObj struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type mapRepo struct {
	URLs       map[string]string
	UserUrls   map[string][]URLObj
	PostgresDB *pgx.Conn
}

func NewSomeRepo(postgres *PostgresRepo) *mapRepo {
	return &mapRepo{
		URLs:       make(map[string]string),
		UserUrls:   make(map[string][]URLObj),
		PostgresDB: postgres.Conn,
	}
}

func (r *mapRepo) Add(url string) (string, error) {
	for {
		uuid := GenUUID()
		if _, ok := r.URLs[uuid]; !ok {
			r.URLs[uuid] = url
			return uuid, nil
		}
	}
}

func (r *mapRepo) Get(uuid string) (string, error) {
	if uuid == "" {
		return "", errors.New("id of url isn't set")
	}

	url := r.URLs[uuid]

	if url == "" {
		return "", errors.New("url with that id is not found")
	}

	return url, nil
}

func (r *mapRepo) AddByUser(userID, originalURL, shortURL string) {

	r.UserUrls[userID] = append(r.UserUrls[userID], URLObj{OriginalURL: originalURL, ShortURL: shortURL})

}

func (r *mapRepo) AllUsersURLS(userID string) []URLObj {
	return r.UserUrls[userID]
}

func (r *mapRepo) IsUserExist(userID string) bool {
	_, ok := r.UserUrls[userID]
	return ok
}

func (r *mapRepo) Ping() error {
	ctx := context.Background()
	if r.PostgresDB == nil {
		return errors.New("something wrong with DB connection")
	}
	return r.PostgresDB.Ping(ctx)
}
