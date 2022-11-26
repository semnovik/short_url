package handlers

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func (h *Handler) GetFullURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")

	URL, err := h.Service.GetURLByID(urlID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) SendURL(w http.ResponseWriter, r *http.Request) {
	request, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urlID := h.Service.PostURL(string(request))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + urlID))
}
