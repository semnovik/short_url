package service

import (
	"short_url/internal/app/repository"
)

type ShorterService interface {
	PostURL(url string) string
	GetURLByID(urlID string) (string, error)
}

type Service struct {
	ShorterService
}

func NewServer(repos *repository.Repository) *Service {
	return &Service{
		ShorterService: NewShorter(repos),
	}
}