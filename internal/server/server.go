package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"short_url/internal/services"
)

type srv struct {
	Service services.Shorter
}

func New(service services.Shorter) *http.Server {
	handler := &srv{Service: service}

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/{id}", handler.GetFullURL)
	router.Post("/", handler.SendURL)

	return &http.Server{Handler: router, Addr: viper.GetString("app.port")}
}

func (h *srv) GetFullURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")

	URL, err := h.Service.GetURLByID(urlID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *srv) SendURL(w http.ResponseWriter, r *http.Request) {
	request, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urlID := h.Service.PostURL(string(request))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + urlID))
}
