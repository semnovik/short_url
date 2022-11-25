package handlers

import (
	"io"
	"net/http"
	"strings"
)

func (h *Handler) startPage(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		urlID := strings.Trim(r.URL.Path, "/")

		URL, err := h.Service.GetURLByID(urlID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", URL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case r.Method == http.MethodPost:
		defer r.Body.Close()

		request, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		urlID := h.Service.PostURL(string(request))

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + urlID))

	default:
		http.Error(w, `Method not found`, http.StatusBadRequest)
	}
}
