package main

import (
	"log"
	"net/http"
	"short_url/internal/handlers"
)

func main() {
	err := http.ListenAndServe(":8080", handlers.InitRouter())
	if err != nil {
		log.Fatal(err)
	}
}
