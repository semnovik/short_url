package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/{id}", h.GetFullURL)
	router.Post("/", h.SendURL)

	return router
}
