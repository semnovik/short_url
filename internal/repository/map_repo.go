package repository

import (
	"database/sql"
	"errors"
	"short_url/configs"
)

type URLObj struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Unit struct {
	OriginalURL string `json:"original_url"`
	ShortUUID   string `json:"short_uuid"`
	UserUUID    string `json:"user_uuid"`
	IsDeleted   bool   `json:"is_deleted"`
}

type MemoryRepo struct {
	Units      []*Unit
	PostgresDB *sql.DB
}

type Some struct {
}

func NewSomeRepo() *MemoryRepo {
	return &MemoryRepo{
		Units:      []*Unit{},
		PostgresDB: nil,
	}
}

func (r *MemoryRepo) Get(uuid string) (string, bool, error) {
	if uuid == "" {
		return "", false, errors.New("id of url isn't set")
	}

	var url string
	for _, obj := range r.Units {
		if uuid == obj.ShortUUID {
			if obj.IsDeleted {
				return "", true, nil
			}
			url = obj.OriginalURL

		}
	}

	if url == "" {
		return "", false, errors.New("url with that id is not found")
	}

	return url, false, nil
}

func (r *MemoryRepo) AllUsersURLS(userID string) []URLObj {
	var result []URLObj

	for _, obj := range r.Units {
		if userID == obj.UserUUID {
			result = append(result, URLObj{OriginalURL: obj.OriginalURL, ShortURL: configs.Config.BaseURL + "/" + obj.ShortUUID})
		}
	}

	return result
}

func (r *MemoryRepo) IsUserExist(userUUID string) bool {
	for _, obj := range r.Units {
		if userUUID == obj.UserUUID {
			return true
		}
	}

	return false
}

func (r *MemoryRepo) Ping() error {
	if r.PostgresDB == nil {
		return errors.New("something wrong with DB connection")
	}
	return r.PostgresDB.Ping()
}

func (r *MemoryRepo) AddByUser(userID, originalURL string) (string, error) {
	var uuid string

	for _, obj := range r.Units {
		if originalURL == obj.OriginalURL {
			return obj.ShortUUID, errors.New(`already exists`)
		}
	}

	uuid = GenUUID()
	r.Units = append(r.Units, &Unit{
		OriginalURL: originalURL,
		ShortUUID:   uuid,
		UserUUID:    userID,
		IsDeleted:   false,
	})

	return uuid, nil
}

func (r *MemoryRepo) DeleteByUUID(uuids []string, userID string) {
	for _, uuid := range uuids {
		for _, obj := range r.Units {
			if userID == obj.UserUUID && uuid == obj.ShortUUID && !obj.IsDeleted {
				obj.IsDeleted = true
			}
		}
	}
}
