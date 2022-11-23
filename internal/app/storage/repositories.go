package storage

import "errors"

type URLRepo struct {
	URLs []string
}

func NewURLRepo() *URLRepo {
	return &URLRepo{URLs: []string{}}
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
