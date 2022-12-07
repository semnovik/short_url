package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"short_url/internal/app/services"
)

type Handler struct {
	Service services.ShorterService
}

func NewHandler(service services.ShorterService) http.Handler {
	handler := &Handler{Service: service}

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/{id}", handler.GetFullURL)
	router.Post("/", handler.SendURL)

	return router
}
