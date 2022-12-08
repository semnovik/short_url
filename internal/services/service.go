package services

import (
	"errors"
	"math/rand"
	"short_url/internal/repositories"
	"time"
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
	for {
		uuid := generateUUID()

		_, err := s.Repository.Get(uuid)
		// Error returns only when url not founded by uuid
		if err != nil {
			s.Repository.Add(uuid, url)
			return uuid
		}
	}

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

func generateUUID() string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	randUUID := make([]byte, 5)

	for i := range randUUID {
		randUUID[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(randUUID)
}
