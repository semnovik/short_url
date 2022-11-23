package server

import (
	"errors"
	"log"
	"short_url/internal/app/storage"
	"strconv"
)

func PostURL(url string) string {

	log.Print("got PostURL request: " + url)

	urlID := strconv.Itoa(storage.Counter + 1)
	storage.UrlsMap[urlID] = url

	return urlID
}

func GetURLByID(urlID string) (string, error) {
	log.Printf("got GetURLByID request: " + urlID)

	if urlID == "" {
		return "", errors.New("id of url isn't set")
	}

	URL, ok := storage.UrlsMap[urlID]
	if !ok {
		return "", errors.New("url with id:" + urlID + " is not found")
	}

	return URL, nil
}
