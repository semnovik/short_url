package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"os"
	"short_url/configs"
	"time"
)

//go:generate mockgen -source=url_repo.go -destination=mock/mock.go

type URLRepo interface {
	Add(url string) (string, error)
	Get(uuid string) (url string, err error)
}

type fileRepo struct {
	mapRepo mapRepo
	File    *os.File
}

type mapRepo struct {
	URLs map[string]string
}

func NewURLRepository() (URLRepo, error) {
	if configs.Config.FileStoragePath != "" {
		file, urls, err := fillRepoFromFile()
		if err != nil {
			return nil, err
		}

		return &fileRepo{
			mapRepo: mapRepo{urls},
			File:    file,
		}, nil
	}

	return &mapRepo{URLs: make(map[string]string)}, nil
}

func (r *mapRepo) Add(url string) (string, error) {
	for {
		uuid := genUUID()
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

func (r *fileRepo) Add(url string) (string, error) {
	for {
		uuid := genUUID()
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

func (r *fileRepo) Get(uuid string) (string, error) {
	return r.mapRepo.Get(uuid)
}

var genUUID = generateUUID

func generateUUID() string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	randUUID := make([]byte, 5)

	for i := range randUUID {
		randUUID[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(randUUID)
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
