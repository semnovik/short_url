package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"os"
	"short_url/configs"
	"time"
)

//go:generate mockgen -source=url_repo.go -destination=mock/mock.go

type URLRepo interface {
	Add(url string) string
	Get(uuid string) (url string, err error)
}

type repoURL struct {
	URLs map[string]string
	File *os.File
}

func NewURLRepository() *repoURL {
	if configs.Config.FileStoragePath != "" {
		file, urls := fillRepoFromFile()

		return &repoURL{
			URLs: urls,
			File: file,
		}
	}

	return &repoURL{
		URLs: make(map[string]string),
		File: nil,
	}
}

func (r *repoURL) Add(url string) string {
	for {
		uuid := genUUID()
		if _, ok := r.URLs[uuid]; !ok {

			// Если есть интеграция с файлом, то пишем еще и в файл
			if r.File != nil {
				writeURLInFile(r.File, uuid, url)
			}

			r.URLs[uuid] = url
			return uuid
		}
	}
}

func (r *repoURL) Get(uuid string) (string, error) {
	if uuid == "" {
		return "", errors.New("id of url isn't set")
	}

	url := r.URLs[uuid]

	if url == "" {
		return "", errors.New("url with that id is not found")
	}

	return url, nil
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

func writeURLInFile(file *os.File, uuid string, url string) {
	writer := bufio.NewWriter(file)
	event := Event{
		UUID: uuid,
		URL:  url,
	}
	data, err := json.Marshal(event)
	if err != nil {
		log.Print(err)
	}

	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}

	err = writer.WriteByte('\n')
	if err != nil {
		log.Print(err)
	}

	err = writer.Flush()
	if err != nil {
		log.Print(err)
	}
}

func fillRepoFromFile() (*os.File, map[string]string) {
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
			log.Print(err)
		}

		urls[event.UUID] = event.URL
	}

	return file, urls
}
