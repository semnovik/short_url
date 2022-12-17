package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"short_url/configs"
	"short_url/internal/repository"
)

type shorterSrv struct {
	repo repository.URLRepo
}

func NewShorterSrv(repo repository.URLRepo) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	h := &shorterSrv{repo: repo}
	router.Post("/api/shorten", h.Shorten)
	router.Get("/{id}", h.GetFullURL)
	router.Post("/", h.SendURL)

	return &http.Server{Handler: router, Addr: configs.Config.ServerAddress}
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

	urlID, err := h.repo.Add(string(request))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(configs.Config.BaseURL + "/" + urlID))
}

type RequestShorten struct {
	URL string `json:"url"`
}
type ResponseShorten struct {
	Result string `json:"result"`
}

func (h *shorterSrv) Shorten(w http.ResponseWriter, r *http.Request) {

	req := RequestShorten{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(data, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	uuid, err := h.repo.Add(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	shortenURL := configs.Config.BaseURL + "/" + uuid

	respBody := ResponseShorten{Result: shortenURL}
	response, _ := json.Marshal(respBody)

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
