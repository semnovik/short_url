package repository

import "errors"

type RepoURL struct {
	URLs []string
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
