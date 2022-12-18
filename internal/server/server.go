package server

import (
	"compress/gzip"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"short_url/configs"
	"short_url/internal/repository"
	"strings"
)

type shorterSrv struct {
	repo repository.URLRepo
}

func NewShorterSrv(repo repository.URLRepo) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))
	router.Use(gzipHandle)
	router.Use(gzipSend)

	h := &shorterSrv{repo: repo}
	router.Post("/api/shorten", h.Shorten)
	router.Get("/{id}", h.GetFullURL)
	router.Post("/", h.SendURL)

	return &http.Server{Handler: router, Addr: configs.Config.ServerAddress}
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func gzipSend(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(writer, request)
			return
		}
		gz, err := gzip.NewWriterLevel(writer, gzip.BestSpeed)
		if err != nil {
			io.WriteString(writer, err.Error())
			return
		}
		defer gz.Close()

		writer.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{ResponseWriter: writer, Writer: gz}, request)
	})
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(request.Body)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			request.Body = gz
		}

		next.ServeHTTP(writer, request)
	})
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
