package server

import (
	"errors"
	"log"
	"short_url/internal/app/storage"
	"strconv"
)

func PostURL(url string) int {

	log.Print("got PostURL request: " + url)

	storage.UrlRepo = append(storage.UrlRepo, url)

	urlID := len(storage.UrlRepo)

	return urlID
}

func GetURLByID(urlID string) (string, error) {
	log.Printf("got GetURLByID request: " + urlID)

	if urlID == "" {
		return "", errors.New("id of url isn't set")
	}

	id, err := strconv.Atoi(urlID)
	if err != nil {
		return "", errors.New("something went wrong")
	}

	if id > len(storage.UrlRepo) || id <= 0 {
		return "", errors.New("url with id:" + urlID + " is not found")
	}

	URL := storage.UrlRepo[id-1]

	return URL, nil
}
