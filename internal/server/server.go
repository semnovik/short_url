package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"short_url/internal/repository"
)

type shorterSrv struct {
	repo repository.URLRepo
}

func NewShorterSrv(repo repository.URLRepo) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	h := &shorterSrv{repo: repo}
	router.Get("/{id}", h.GetFullURL)
	router.Post("/", h.SendURL)

	return &http.Server{Handler: router, Addr: viper.GetString("app.port")}
}

func (h *shorterSrv) GetFullURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")

	URL, err := h.repo.Get(urlID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *shorterSrv) SendURL(w http.ResponseWriter, r *http.Request) {
	request, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urlID := h.repo.Add(string(request))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + urlID))
}
