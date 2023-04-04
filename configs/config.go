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
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"./file_with_urls"`
	DatabaseDSN     string `env:"DATABASE_DSN"` // host=localhost port=5438 dbname=admin user=admin password=password
	MigrationsDir   string `env:"MIGRATIONS_DIR" envDefault:"./migrations"`
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
		flag.StringVar(&Config.FileStoragePath, "f", "./file_with_urls", "path to file with urls")
	}
	if _, exist := os.LookupEnv("DATABASE_DSN"); !exist {
		flag.StringVar(&Config.DatabaseDSN, "d", "host=localhost port=5438 dbname=admin user=admin password=password", "dsn for db")
		// host=localhost port=5438 dbname=admin user=admin password=password - local
		// host=db port=5432 dbname=admin user=admin password=password - inside docker
	}
	if _, exist := os.LookupEnv("MIGRATIONS_DIR"); !exist {
		flag.StringVar(&Config.MigrationsDir, "mg", "./migrations", "path to migration files")
	}
	flag.Parse()
}
