package server

import (
	"errors"
	"log"
	"short_url/internal/app/storage"
	"strconv"
)

type Server struct {
	Repository *storage.Repository
}

func NewServer(repos *storage.Repository) *Server {
	return &Server{
		Repository: repos,
	}
}

func (s *Server) PostURL(url string) string {

	log.Print("got PostURL request: " + url)

	urlID := strconv.Itoa(s.Repository.AddURL(url))

	return urlID
}

func (s *Server) GetURLByID(urlID string) (string, error) {
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
