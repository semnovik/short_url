package handlers

import (
	"io"
	"net/http"
	"strings"
)

func (h *Handler) startPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	switch {
	case r.Method == http.MethodGet:
		urlID := strings.Trim(r.URL.Path, "/")

		URL, err := h.Server.GetURLByID(urlID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", URL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case r.Method == http.MethodPost:
		request, err := io.ReadAll(r.Body)
		if err != nil || r.Body == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		urlID := h.Server.PostURL(string(request))

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + urlID))

	default:
		http.Error(w, "Method not found", http.StatusBadRequest)
	}
}
