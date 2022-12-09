package repository

import (
	"errors"
	"math/rand"
	"time"
)

//go:generate mockgen -source=url_repo.go -destination=mock/mock.go

type URLRepo interface {
	Add(url string) string
	Get(uuid string) (url string, err error)
}

type repoURL struct {
	URLs map[string]string
}

func NewURLRepository() *repoURL {
	return &repoURL{URLs: make(map[string]string)}
}

func (r *repoURL) Add(url string) string {
	for {
		uuid := genUUID()

		if _, ok := r.URLs[uuid]; !ok {
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
