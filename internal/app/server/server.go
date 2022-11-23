package server

import (
	"errors"
	"log"
	"short_url/internal/app/storage"
	"strconv"
)

type ShorterService interface {
	PostURL(url string) string
	GetURLByID(urlID string) (string, error)
}

type Shorter struct {
	Repository *storage.Repository
}

type Server struct {
	ShorterService
}

func NewServer(repos *storage.Repository) *Server {
	return &Server{
		ShorterService: NewShorter(repos),
	}
}

func NewShorter(repos *storage.Repository) *Shorter {
	return &Shorter{
		Repository: repos,
	}
}

func (s *Shorter) PostURL(url string) string {

	log.Print("got PostURL request: " + url)

	urlID := strconv.Itoa(s.Repository.AddURL(url))

	return urlID
}

func (s *Shorter) GetURLByID(urlID string) (string, error) {
	log.Printf("got GetURLByID request: " + urlID)

	if urlID == "" {
		return "", errors.New("id of url isn't set")
	}

	id, err := strconv.Atoi(urlID)
	if err != nil {
		return "", errors.New("something went wrong")
	}

	URL, err := s.Repository.GetURL(id)
	if err != nil {
		return "", err
	}

	return URL, nil
}
