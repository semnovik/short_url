package repository

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
)

type FileRepo struct {
	mapRepo    MapRepo
	File       *os.File
	PostgresDB *sql.DB
}

func NewFileRepo() *FileRepo {
	file, urls, err := fillRepoFromFile()
	if err != nil {
		return nil
	}

	return &FileRepo{
		mapRepo:    MapRepo{URLs: urls, UserUrls: make(map[string][]URLObj)},
		File:       file,
		PostgresDB: nil,
	}
}

func (r *FileRepo) Add(url string) (string, error) {
	for {
		uuid := GenUUID()
		if _, ok := r.mapRepo.URLs[uuid]; !ok {

			// Если есть интеграция с файлом, то пишем еще и в файл
			if r.File != nil {
				err := writeURLInFile(r.File, uuid, url)
				if err != nil {
					return "", err
				}
			}

			r.mapRepo.URLs[uuid] = url
			return uuid, nil
		}
	}
}

func (r *FileRepo) Get(uuid string) (string, error) {
	return r.mapRepo.Get(uuid)
}

//func (r *FileRepo) AddByUser(userID, originalURL, shortURL string) {
//	r.mapRepo.UserUrls[userID] = append(r.mapRepo.UserUrls[userID], URLObj{OriginalURL: originalURL, ShortURL: shortURL})
//}

func (r *FileRepo) AllUsersURLS(userID string) []URLObj {
	return r.mapRepo.UserUrls[userID]
}

func (r *FileRepo) IsUserExist(userID string) bool {
	_, ok := r.mapRepo.UserUrls[userID]
	return ok
}

func (r *FileRepo) Ping() error {
	if r.PostgresDB == nil {
		return errors.New("something wrong with DB-connection")
	}
	return r.PostgresDB.Ping()
}

func (r *FileRepo) AddByUser(userID, originalURL string) (string, error) {
	var uuid string

	for {
		uuid = GenUUID()
	}
	return uuid, nil
}

type Event struct {
	UUID string `json:"UUID"`
	URL  string `json:"URL"`
}

func writeURLInFile(file *os.File, uuid string, url string) error {
	writer := bufio.NewWriter(file)
	event := Event{
		UUID: uuid,
		URL:  url,
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	err = writer.WriteByte('\n')
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
