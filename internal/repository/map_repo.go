package repository

import (
	"database/sql"
	"errors"
)

type URLObj struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type MapRepo struct {
	URLs       map[string]string
	UserUrls   map[string][]URLObj
	PostgresDB *sql.DB
}

func NewSomeRepo(postgres *PostgresRepo) *MapRepo {
	return &MapRepo{
		URLs:       make(map[string]string),
		UserUrls:   make(map[string][]URLObj),
		PostgresDB: postgres.Conn,
	}
}

func (r *MapRepo) Add(url string) (string, error) {
	for {
		uuid := GenUUID()
		if _, ok := r.URLs[uuid]; !ok {
			r.URLs[uuid] = url
			return uuid, nil
		}
	}
}

func (r *MapRepo) Get(uuid string) (string, error) {
	if uuid == "" {
		return "", errors.New("id of url isn't set")
	}

	url := r.URLs[uuid]

	if url == "" {
		return "", errors.New("url with that id is not found")
	}

	return url, nil
}

func (r *MapRepo) AddByUser(userID, originalURL, shortURL string) {

	r.UserUrls[userID] = append(r.UserUrls[userID], URLObj{OriginalURL: originalURL, ShortURL: shortURL})

}

func (r *MapRepo) AllUsersURLS(userID string) []URLObj {
	return r.UserUrls[userID]
}

func (r *MapRepo) IsUserExist(userID string) bool {
	_, ok := r.UserUrls[userID]
	return ok
}

func (r *MapRepo) Ping() error {
	if r.PostgresDB == nil {
		return errors.New("something wrong with DB connection")
	}
	return r.PostgresDB.Ping()
}
