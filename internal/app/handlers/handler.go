package handlers

import (
	"net/http"
	"short_url/internal/app/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) InitRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", h.startPage)

	return router
}
