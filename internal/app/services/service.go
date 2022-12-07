package services

import (
	"errors"
	"short_url/internal/app/repositories"
)

//go:generate mockgen -source=service.go -destination=mock_service/mock.go

type ShorterService interface {
	PostURL(url string) string
	GetURLByID(urlID string) (string, error)
}

type Shorter struct {
	Repository repositories.URLRepo
}

func NewShorter(repo repositories.URLRepo) *Shorter {
	return &Shorter{
		Repository: repo,
	}
}

func (s *Shorter) PostURL(url string) string {

	uuid := s.Repository.Add(url)

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
