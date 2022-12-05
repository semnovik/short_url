package main

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	"short_url/internal/app/handlers"
	"short_url/internal/app/repositories"
	"short_url/internal/app/services"
)

func main() {

	if err := InitConfig(); err != nil {
		log.Fatal("error with reading config", err)
	}

	repository := repositories.NewRepository()
	shorter := services.NewShorter(repository)
	service := services.NewService(shorter)
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
