package main

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	"short_url/internal/handlers"
	"short_url/internal/repositories"
	"short_url/internal/services"
)

func main() {

	if err := InitConfig(); err != nil {
		log.Fatal("error with reading config", err)
	}

	repository := repositories.NewURLRepo()
	service := services.NewShorter(repository)
	handler := handlers.NewHandler(service)

	err := http.ListenAndServe(viper.GetString("app.port"), handler)
	if err != nil {
		log.Fatal(err)
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
