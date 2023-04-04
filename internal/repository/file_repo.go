package repository

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"os"
	"short_url/configs"
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

func (r *FileRepo) Get(uuid string) (string, bool, error) {
	return r.mapRepo.Get(uuid)
}

func (r *FileRepo) AllUsersURLS(userID string) []URLObj {
	var result []URLObj
	for _, part := range r.mapRepo.UserUrls[userID] {
		part.ShortURL = configs.Config.BaseURL + "/" + part.ShortURL
		result = append(result, part)
	}
	return result
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

	for uuidMemo, origFromMemo := range r.mapRepo.URLs {
		if origFromMemo == originalURL {
			return uuidMemo, errors.New(`already exists`)
		}
	}

	for {
		uuid = GenUUID()
		if _, ok := r.mapRepo.URLs[uuid]; !ok {
			r.mapRepo.UserUrls[userID] = append(r.mapRepo.UserUrls[userID], URLObj{OriginalURL: originalURL, ShortURL: uuid})
			r.mapRepo.URLs[uuid] = originalURL

			if r.File != nil {
				err := writeURLInFile(r.File, uuid, originalURL)
				if err != nil {
					return "", err
				}
			}
			break
		}
	}

	return uuid, nil
}

func (r *FileRepo) DeleteByUUID(uuid []string, userID string) {
	r.mapRepo.DeleteByUUID(uuid, userID)
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

func fillRepoFromFile() (*os.File, map[string]string, error) {
	file, err := os.OpenFile(configs.Config.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		file = nil
	}
	urls := make(map[string]string)

	reader := bufio.NewReader(file)
	for {
		data, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		}

		event := Event{}
		err = json.Unmarshal(data, &event)
		if err != nil {
			return nil, nil, err
		}

		urls[event.UUID] = event.URL
	}

	return file, urls, nil
}
