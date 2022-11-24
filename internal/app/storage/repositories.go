package storage

import "errors"

type URLRepo interface {
	Add(url string) (urlID int)
	Get(urlID int) (url string, err error)
}

type RepoURL struct {
	URLs []string
}

type Repository struct {
	URLRepo
}

func NewRepository(storage []string) *Repository {
	return &Repository{
		URLRepo: NewURLRepo(storage),
	}
}

func NewURLRepo(storage []string) *RepoURL {
	return &RepoURL{URLs: storage}
}

func (r *RepoURL) Add(url string) int {
	r.URLs = append(r.URLs, url)
	urlID := len(r.URLs)

	return urlID
}

func (r *RepoURL) Get(id int) (string, error) {
	if id > len(r.URLs) || id <= 0 {
		return "", errors.New("url with that id is not found")
	}

	return r.URLs[id-1], nil
}
