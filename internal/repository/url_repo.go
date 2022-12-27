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

type URLObj struct {
	originalURL string
	shortURL    string
}

type SomeRepo struct {
	URLs     map[string]string
	UserUrls map[string][]URLObj
}

func NewURLRepository() URLRepo {
	if configs.Config.FileStoragePath != "" {
		return NewFileRepo()
	}

	return newMapRepo()
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
