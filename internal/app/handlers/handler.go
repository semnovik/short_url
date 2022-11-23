package handlers

import (
	"net/http"
	"short_url/internal/app/server"
)

type Handler struct {
	Server *server.Server
}

func NewHandler(server *server.Server) *Handler {
	return &Handler{Server: server}
}

func (h *Handler) InitRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", h.startPage)

	return router
}
