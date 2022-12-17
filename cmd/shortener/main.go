package main

import (
	"flag"
	"log"
	"os"
	"short_url/configs"
	"short_url/internal/repository"
	"short_url/internal/server"
)

func main() {
	initFlags()
	repo, err := repository.NewURLRepository()
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewShorterSrv(repo)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func initFlags() {
	if _, exist := os.LookupEnv("SERVER_ADDRESS"); !exist {
		flag.StringVar(&configs.Config.ServerAddress, "a", ":8080", "address of the server")
	}
	if _, exist := os.LookupEnv("BASE_URL"); !exist {
		flag.StringVar(&configs.Config.BaseURL, "b", "http://localhost:8080", "base URL of the server")
	}
	if _, exist := os.LookupEnv("FILE_STORAGE_PATH"); !exist {
		flag.StringVar(&configs.Config.FileStoragePath, "f", "./internal/repository/file_with_urls", "path to file with urls")
	}

	flag.Parse()
}
