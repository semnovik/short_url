package repository

import (
	"database/sql"
	"errors"
	"short_url/configs"
)

type URLObj struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	IsDeleted   bool   `json:"is_deleted"`
}

type MapRepo struct {
	URLs       map[string]string
	UserUrls   map[string][]URLObj
	PostgresDB *sql.DB
}

func NewSomeRepo() *MapRepo {
	return &MapRepo{
		URLs:       make(map[string]string),
		UserUrls:   make(map[string][]URLObj),
		PostgresDB: nil,
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

func (r *MapRepo) Get(uuid string) (string, bool, error) {
	if uuid == "" {
		return "", false, errors.New("id of url isn't set")
	}

	url := r.URLs[uuid]

	if url == "" {
		return "", false, errors.New("url with that id is not found")
	}

	return url, false, nil
}

func (r *MapRepo) AllUsersURLS(userID string) []URLObj {
	var result []URLObj
	for _, part := range r.UserUrls[userID] {
		part.ShortURL = configs.Config.BaseURL + "/" + part.ShortURL
		result = append(result, part)
	}
	return result
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

func (r *MapRepo) AddByUser(userID, originalURL string) (string, error) {
	var uuid string

	for uuidMemo, origFromMemo := range r.URLs {
		if origFromMemo == originalURL {
			return uuidMemo, errors.New(`already exists`)
		}
	}

	for {
		uuid = GenUUID()
		if _, ok := r.URLs[uuid]; !ok {
			r.UserUrls[userID] = append(r.UserUrls[userID], URLObj{OriginalURL: originalURL, ShortURL: uuid})
			r.URLs[uuid] = originalURL
			break
		}
	}

	return uuid, nil
}

func (r *MapRepo) DeleteByUUID(uuid, userID string) {
	uuid = "заглушка"
	userID = "заглушка"
}
