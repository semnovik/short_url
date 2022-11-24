package handlers

import (
	"net/http"
	"short_url/internal/app/service"
)

type Handler struct {
	Server *service.Service
}

func NewHandler(server *service.Service) *Handler {
	return &Handler{Server: server}
}

func (h *Handler) InitRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", h.startPage)

	return router
}
