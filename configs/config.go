package configs

import (
	"github.com/caarlos0/env/v6"
	"log"
)

var Config = InitConfig()

type Conf struct {
	ServerAddress string `env:"ADDRESS,required"`
	BaseURL       string `env:"BASE_URL,required"`
}

func InitConfig() Conf {
	var cfg Conf

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
