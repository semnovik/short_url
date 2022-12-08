package repositories

import (
	"errors"
)

//go:generate mockgen -source=repositories.go -destination=mock_repositories/mock.go

type URLRepo interface {
	Add(id string, url string)
	Get(id string) (url string, err error)
}

var urlStorage = make(map[string]string)

type repoURL struct {
	URLs map[string]string
}

func NewRepository() *repoURL {
	return &repoURL{URLs: urlStorage}
}

func (r *repoURL) Add(id, url string) {
	r.URLs[id] = url
}

func (r *repoURL) Get(uuid string) (string, error) {
	url := r.URLs[uuid]

	if url == "" {
		return "", errors.New("url with that id is not found")
	}

	return url, nil
}
