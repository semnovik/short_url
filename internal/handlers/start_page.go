package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"short_url/internal/server"
	"strings"
)

func startPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	switch {
	case r.Method == http.MethodGet:
		urlID := strings.Trim(r.URL.Path, "/")

		URL, err := server.GetURLByID(urlID)

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

		urlID := server.PostURL(string(request))

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + urlID))

	default:
		NotFound(w, r)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {

	msg := InfoMessage{Message: "Method not found"}
	res, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
}
