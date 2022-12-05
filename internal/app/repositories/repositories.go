package repositories

import (
	"errors"
	"math/rand"
	"time"
)

type URLRepo interface {
	Add(url string) (uuid string)
	Get(urlID string) (url string, err error)
}

var urlStorage = make(map[string]string)

type RepoURL struct {
	URLs map[string]string
}

func NewURLRepo() *RepoURL {
	return &RepoURL{URLs: urlStorage}
}

func (r *RepoURL) Add(url string) string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	randUUID := make([]byte, 5)
	for i := range randUUID {
		randUUID[i] = charset[seededRand.Intn(len(charset))]
	}

	r.URLs[string(randUUID)] = url

	return string(randUUID)
}

func (r *RepoURL) Get(uuid string) (string, error) {
	url := r.URLs[uuid]

	if url == "" {
		{
			return "", errors.New("url with that id is not found")
		}
	}

	return url, nil
}
