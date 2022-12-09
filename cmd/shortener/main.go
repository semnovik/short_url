package main

import (
	"github.com/spf13/viper"
	"log"
	"short_url/internal/repository"
	"short_url/internal/server"
)

func main() {
	if err := InitConfig(); err != nil {
		log.Fatal("error with reading config", err)
	}

	repository := repository.NewURLRepository()
	srv := server.New(repository)

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
