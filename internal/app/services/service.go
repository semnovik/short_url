package services

import (
	"short_url/internal/app/repositories"
)

type ShorterService interface {
	PostURL(url string) string
	GetURLByID(urlID string) (string, error)
}

type Service struct {
	ShorterService
}

func NewServer(repos *repositories.Repository) *Service {
	return &Service{
		ShorterService: NewShorter(repos),
	}
}
