package repository

import (
	"database/sql"
	"log"
	"math/rand"
	"short_url/configs"
	"time"
)

//go:generate mockgen -source=repos.go -destination=mock/mock.go

type URLStorage interface {
	Add(url string) (string, error)
	Get(uuid string) (string, bool, error)
	AddByUser(userID, originalURL string) (string, error)
	AllUsersURLS(userID string) []URLObj
	IsUserExist(userID string) bool
	Ping() error
	DeleteByUUID(uuid, userId string)
}

func NewRepo(db *sql.DB) URLStorage {
	if db != nil {
		log.Print("Selected Postgres DB for repository")
		return NewPostgresRepo(db)
	}

	if configs.Config.FileStoragePath != "" {
		log.Print("Selected FileStorage for repository")
		return NewFileRepo()
	}

	log.Print("Selected InMemory for repository")
	return NewSomeRepo()
}

var GenUUID = generateUUID

func generateUUID() string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	randUUID := make([]byte, 10)

	for i := range randUUID {
		randUUID[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(randUUID)
}
