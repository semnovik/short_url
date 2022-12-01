package main

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	"short_url/internal/app/handlers"
	"short_url/internal/app/repository"
	"short_url/internal/app/service"
)

func main() {

	if err := InitConfig(); err != nil {
		log.Fatal("error with reading config", err)
	}

	Repository := repository.NewRepository()
	Server := service.NewServer(Repository)
	Handler := handlers.NewHandler(Server)

	err := http.ListenAndServe(viper.GetString("app.port"), Handler.InitRouter())
	if err != nil {
		log.Fatal(err)
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
