package repositories

import (
	"errors"
	"math/rand"
	"time"
)

type RepoURL struct {
	URLs map[string]string
}

func NewURLRepo(storage map[string]string) *RepoURL {
	return &RepoURL{URLs: storage}
}

func (r *RepoURL) Add(url string) string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	randUuid := make([]byte, 5)
	for i := range randUuid {
		randUuid[i] = charset[seededRand.Intn(len(charset))]
	}

	r.URLs[string(randUuid)] = url

	return string(randUuid)
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
