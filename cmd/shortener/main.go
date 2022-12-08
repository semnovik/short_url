package main

import (
	"github.com/spf13/viper"
	"log"
	"short_url/internal/server"
	"short_url/internal/service"
	"short_url/internal/storage"
)

func main() {
	if err := InitConfig(); err != nil {
		log.Fatal("error with reading config", err)
	}

	repository := storage.NewURLRepository()
	shorterService := service.NewShorter(repository)
	srv := server.New(shorterService)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
