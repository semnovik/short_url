package services

import (
	"errors"
	"short_url/internal/app/repositories"
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

	uuid := s.Repository.URLRepo.Add(url)

	return uuid
}

func (s *Shorter) GetURLByID(uuid string) (string, error) {

	if uuid == "" {
		return "", errors.New("id of url isn't set")
	}

	URL, err := s.Repository.Get(uuid)
	if err != nil {
		return "", err
	}

	return URL, nil
}
