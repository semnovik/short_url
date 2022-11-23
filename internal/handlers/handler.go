package handlers

import (
	"net/http"
)

type InfoMessage struct {
	Message string `json:"message"`
}

func InitRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", startPage)

	return router
}
