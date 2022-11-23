package storage

import "errors"

type SliceRepo interface {
	AddURL(url string) int
	GetURL(id int) (string, error)
}

type URLRepo struct {
	URLs []string
}

type Repository struct {
	SliceRepo
}

func NewRepository(storage []string) *Repository {
	return &Repository{
		SliceRepo: NewURLRepo(storage),
	}
}

func NewURLRepo(storage []string) *URLRepo {
	return &URLRepo{URLs: storage}
}

func (r *URLRepo) AddURL(url string) int {
	r.URLs = append(r.URLs, url)
	urlID := len(r.URLs)

	return urlID
}

func (r *URLRepo) GetURL(id int) (string, error) {
	if id > len(r.URLs) || id <= 0 {
		return "", errors.New("url with that id is not found")
	}

	return r.URLs[id-1], nil
}
