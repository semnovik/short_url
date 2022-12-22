package configs

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
)

var Config cfg

type cfg struct {
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

func InitFlags() {
	if _, exist := os.LookupEnv("SERVER_ADDRESS"); !exist {
		flag.StringVar(&Config.ServerAddress, "a", ":8080", "address of the server")
	}
	if _, exist := os.LookupEnv("BASE_URL"); !exist {
		flag.StringVar(&Config.BaseURL, "b", "http://localhost:8080", "base URL of the server")
	}
	if _, exist := os.LookupEnv("FILE_STORAGE_PATH"); !exist {
		flag.StringVar(&Config.FileStoragePath, "f", "./internal/repository/file_with_urls", "path to file with urls")
	}

	flag.Parse()
}
