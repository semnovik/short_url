package repository

import (
	"bufio"
	"encoding/json"
	"os"
)

type fileRepo struct {
	mapRepo mapRepo
	File    *os.File
}

func NewFileRepo() *fileRepo {
	file, urls, err := fillRepoFromFile()
	if err != nil {
		return nil
	}

	return &fileRepo{
		mapRepo: mapRepo{urls},
		File:    file,
	}
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
