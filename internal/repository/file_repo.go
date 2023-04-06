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
	memoryRepo MemoryRepo
	File       *os.File
	PostgresDB *sql.DB
}

func NewFileRepo() *FileRepo {
	file, units, err := fillRepoFromFile()
	if err != nil {
		return nil
	}

	return &FileRepo{
		memoryRepo: MemoryRepo{Units: units},
		File:       file,
		PostgresDB: nil,
	}
}

func (r *FileRepo) Get(uuid string) (string, bool, error) {
	return r.memoryRepo.Get(uuid)
}

func (r *FileRepo) AllUsersURLS(userID string) []URLObj {
	return r.memoryRepo.AllUsersURLS(userID)
}

func (r *FileRepo) IsUserExist(userUUID string) bool {
	return r.memoryRepo.IsUserExist(userUUID)
}

func (r *FileRepo) Ping() error {
	return r.memoryRepo.Ping()
}

func (r *FileRepo) AddByUser(userID, originalURL string) (string, error) {
	var uuid string

	for _, obj := range r.memoryRepo.Units {
		if originalURL == obj.OriginalURL {
			return obj.ShortUUID, errors.New(`already exists`)
		}
	}

	uuid = GenUUID()
	r.memoryRepo.Units = append(r.memoryRepo.Units, &Unit{
		OriginalURL: originalURL,
		ShortUUID:   uuid,
		UserUUID:    userID,
		IsDeleted:   false,
	})

	if r.File != nil {
		err := writeURLInFile(r.File, uuid, originalURL, userID)
		if err != nil {
			return "", err
		}
	}

	return uuid, nil
}

func (r *FileRepo) DeleteByUUID(uuid []string, userID string) {
	r.memoryRepo.DeleteByUUID(uuid, userID)
}

func writeURLInFile(file *os.File, uuid string, url string, userUUID string) error {
	writer := bufio.NewWriter(file)
	event := Unit{
		OriginalURL: url,
		ShortUUID:   uuid,
		UserUUID:    userUUID,
		IsDeleted:   false,
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

func fillRepoFromFile() (*os.File, []*Unit, error) {
	file, err := os.OpenFile(configs.Config.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		file = nil
	}

	var units []*Unit
	reader := bufio.NewReader(file)
	for {
		data, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		}

		unit := Unit{}
		err = json.Unmarshal(data, &unit)
		if err != nil {
			return nil, nil, err
		}

		units = append(units, &unit)
	}

	return file, units, nil
}
