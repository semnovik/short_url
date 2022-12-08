package repositories

import (
	"errors"
	"math/rand"
	"time"
)

//go:generate mockgen -source=repositories.go -destination=repo_mocks/mock.go

type URLRepo interface {
	Add(url string) (uuid string)
	Get(uuid string) (url string, err error)
}

var urlStorage = make(map[string]string)

type repoURL struct {
	URLs map[string]string
}

func NewURLRepo() *repoURL {
	return &repoURL{URLs: urlStorage}
}

func (r *repoURL) Add(url string) string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	randUUID := make([]byte, 5)
	for i := range randUUID {
		randUUID[i] = charset[seededRand.Intn(len(charset))]
	}

	r.URLs[string(randUUID)] = url

	return string(randUUID)
}

func (r *repoURL) Get(uuid string) (string, error) {
	url := r.URLs[uuid]

	if url == "" {
		return "", errors.New("url with that id is not found")
	}

	return url, nil
}
