package configs

import (
	"github.com/caarlos0/env/v6"
	"log"
)

var Config = InitConfig()

type Conf struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
}

func InitConfig() Conf {
	var cfg Conf

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
