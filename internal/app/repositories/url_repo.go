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
