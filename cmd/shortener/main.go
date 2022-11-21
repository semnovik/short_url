package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type InfoMessage struct {
	Message string `json:"message"`
}

type PostLongURL struct {
	URL string `json:"url"`
}

var counter int
var UrlsMap = make(map[string]string)

func FirstPage(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		GetURLByID(w, r)
	case r.Method == http.MethodPost:
		PostUrl(w, r)
	default:
		notFound(w, r)
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", FirstPage)
	// конструируем сервер
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func GetURLByID(w http.ResponseWriter, r *http.Request) {
	idRow := strings.Trim(r.URL.Path, "/")
	checkIn, ok := UrlsMap[idRow]
	if !ok {
		w.Write([]byte("there is no url with that id"))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Header().Set("Location", checkIn)
	w.WriteHeader(http.StatusTemporaryRedirect)

	msg := InfoMessage{Message: "Success"}

	resp, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.Write(resp)
}

func PostUrl(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)

	if err != nil {
		log.Printf("error occured while reading the PostUrl request body: %s", err.Error())
		http.Error(w, err.Error(), 500)
	}

	log.Printf("got request with body:" + string(req))

	counter++
	res := strconv.Itoa(counter)

	UrlsMap[res] = string(req)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func notFound(w http.ResponseWriter, r *http.Request) {

	msg := InfoMessage{Message: "Method not found"}
	res, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(res)
}
