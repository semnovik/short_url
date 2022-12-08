package services

import (
	"errors"
	"short_url/internal/repositories"
)

//go:generate mockgen -source=service.go -destination=mock_service/mock.go

type ShorterService interface {
	PostURL(url string) string
	GetURLByID(urlID string) (string, error)
}

type shorter struct {
	Repository repositories.URLRepo
}

func NewShorter(repo repositories.URLRepo) *shorter {
	return &shorter{
		Repository: repo,
	}
}

func (s *shorter) PostURL(url string) string {

	uuid := s.Repository.Add(url)

	return uuid
}

func (s *shorter) GetURLByID(uuid string) (string, error) {

	if uuid == "" {
		return "", errors.New("id of url isn't set")
	}

	URL, err := s.Repository.Get(uuid)
	if err != nil {
		return "", err
	}

	return URL, nil
}
