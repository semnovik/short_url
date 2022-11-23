package handlers

import (
	"net/http"
)

func InitRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", startPage)

	return router
}
