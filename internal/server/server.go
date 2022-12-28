package server

import (
	"encoding/json"
	"errors"
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
	router.Use(gzipReceive)
	router.Use(gzipSend)

	h := &shorterSrv{repo: repo}
	router.Post("/api/shorten", h.Shorten)
	router.Get("/{id}", h.GetFullURL)
	router.Post("/", h.SendURL)
	router.Get("/api/user/urls", h.AllUserURLS)

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

	cookies := r.Cookies()
	var userExist bool
	var userId string
	var newUserToken string
	var err error

	for _, v := range cookies {
		if v.Name == "Auth" {
			userId, err = decodeToken(v.Value)
			if err != nil {
				continue
			}
			userExist = h.repo.IsUserExist(userId)
			break
		}
	}

	if !userExist {
		userId, newUserToken = generateEncodedToken()
		newCookie := &http.Cookie{
			Name:  "Auth",
			Value: newUserToken,
		}
		http.SetCookie(w, newCookie)
	}

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

	h.repo.AddByUser(userId, req.URL, configs.Config.BaseURL+"/"+uuid)

	shortenURL := configs.Config.BaseURL + "/" + uuid

	respBody := ResponseShorten{Result: shortenURL}
	response, _ := json.Marshal(respBody)

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (h *shorterSrv) AllUserURLS(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	var userId string
	var err error
	for _, cookie := range cookies {
		if cookie.Name == "Auth" {
			userId, err = decodeToken(cookie.Value)
			if err != nil {
				continue
			}
		}
	}

	urls := h.repo.AllUsersURLS(userId)
	if len(urls) == 0 {
		http.Error(w, errors.New("not found").Error(), http.StatusNoContent)
		return
	}
	response, _ := json.Marshal(urls)
	w.Write(response)
}
