package configs

import (
	"github.com/caarlos0/env/v6"
	"log"
)

var Config Conf

type Conf struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./internal/repository/file_with_urls"`
}

func init() {
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal(err)
	}
}
