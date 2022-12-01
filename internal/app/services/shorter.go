package services

import (
	"errors"
	"short_url/internal/app/repositories"
	"strconv"
)

type Shorter struct {
	Repository *repositories.Repository
}

func NewShorter(repos *repositories.Repository) *Shorter {
	return &Shorter{
		Repository: repos,
	}
}

func (s *Shorter) PostURL(url string) string {

	urlID := strconv.Itoa(s.Repository.Add(url))

	return urlID
}

func (s *Shorter) GetURLByID(urlID string) (string, error) {

	if urlID == "" {
		return "", errors.New("id of url isn't set")
	}

	id, err := strconv.Atoi(urlID)
	if err != nil {
		return "", errors.New("something went wrong")
	}

	URL, err := s.Repository.Get(id)
	if err != nil {
		return "", err
	}

	return URL, nil
}
