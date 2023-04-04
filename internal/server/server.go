package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"short_url/configs"
	"short_url/internal/repository"
	"short_url/pkg"
)

type shorterSrv struct {
	repo repository.URLStorage
}

func NewShorterSrv(repo repository.URLStorage) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))
	router.Use(gzipReceive)
	router.Use(gzipSend)

	h := &shorterSrv{repo: repo}
	router.Post("/api/shorten", h.Shorten)
	router.Get("/{id}", h.GetFullURL)
	router.Post("/", h.SendURL)
	router.Get("/api/user/urls", h.AllUserURLS)
	router.Get("/ping", h.Ping)
	router.Post("/api/shorten/batch", h.Batch)
	router.Delete("/api/user/urls", h.DeleteURLs)

	return &http.Server{Handler: router, Addr: configs.Config.ServerAddress}
}

func (h *shorterSrv) GetFullURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")

	url, isDeleted, err := h.repo.Get(urlID)

	if isDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *shorterSrv) SendURL(w http.ResponseWriter, r *http.Request) {
	userID, isUserExist := checkUserExist(r, h.repo)

	if !isUserExist {
		userID = setNewUserToken(w)
	}

	request, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uuid, err := h.repo.AddByUser(userID, string(request))
	if err != nil {
		if uuid != "" {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(configs.Config.BaseURL + "/" + uuid))
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(configs.Config.BaseURL + "/" + uuid))
}

func (h *shorterSrv) Shorten(w http.ResponseWriter, r *http.Request) {
	userID, isUserExist := checkUserExist(r, h.repo)

	if !isUserExist {
		userID = setNewUserToken(w)
	}

	req := pkg.RequestShorten{}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(data, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	uuid, err := h.repo.AddByUser(userID, req.URL)
	if err != nil {
		if uuid != "" {
			w.WriteHeader(http.StatusConflict)
			shortenURL := configs.Config.BaseURL + "/" + uuid

			respBody := pkg.ResponseShorten{Result: shortenURL}
			response, _ := json.Marshal(respBody)
			w.Write(response)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortenURL := configs.Config.BaseURL + "/" + uuid

	respBody := pkg.ResponseShorten{Result: shortenURL}
	response, _ := json.Marshal(respBody)

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (h *shorterSrv) AllUserURLS(w http.ResponseWriter, r *http.Request) {
	userID, _ := checkUserExist(r, h.repo)

	urls := h.repo.AllUsersURLS(userID)
	if len(urls) == 0 {
		http.Error(w, errors.New("not found").Error(), http.StatusNoContent)
		return
	}
	response, _ := json.Marshal(urls)
	w.Write(response)
}

func (h *shorterSrv) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.repo.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("everything is OK with DB"))
	}
}

func (h *shorterSrv) Batch(w http.ResponseWriter, r *http.Request) {
	var batch []pkg.RequestShortenBatch
	var urls []pkg.ResponseShortenBatch

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(data, &batch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	for _, part := range batch {
		shortURL, err := h.repo.Add(part.OriginalID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		urls = append(urls, pkg.ResponseShortenBatch{CorrelationID: part.CorrelationID, ShortURL: configs.Config.BaseURL + "/" + shortURL})
	}

	response, _ := json.Marshal(urls)
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (h *shorterSrv) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	userID, _ := checkUserExist(r, h.repo)

	request, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var urls []string
	err = json.Unmarshal(request, &urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, uuid := range urls {
		h.repo.DeleteByUUID(uuid, userID)
	}
	fmt.Print(userID)
	w.WriteHeader(http.StatusAccepted)
}
