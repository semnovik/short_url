package main

import (
	"github.com/spf13/viper"
	"log"
	"short_url/internal/repositories"
	"short_url/internal/server"
	"short_url/internal/services"
)

func main() {
	if err := InitConfig(); err != nil {
		log.Fatal("error with reading config", err)
	}

	repository := repositories.NewRepository()
	service := services.NewShorter(repository)

	srv := server.New(service)
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
