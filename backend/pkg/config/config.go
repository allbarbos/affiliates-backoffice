package config

import (
	"log"

	"github.com/caarlos0/env/v7"
)

var Vars appConfig

type appConfig struct {
	ApiLogLevel    string `env:"LOG_LEVEL"`
	ApiTimeout     string `env:"API_TIMEOUT,required"`
	ApiPort        string `env:"API_PORT,required"`
	DbMaxIdleConns int    `env:"DB_MAX_IDLE_CONNS,required"`
	DbMaxOpenConns int    `env:"DB_MAX_OPEN_CONNS,required"`
	DbPort         int    `env:"DB_PORT,required"`
	DbHost         string `env:"DB_HOST,required"`
	DbUser         string `env:"DB_USER,required"`
	DbPassword     string `env:"DB_PASSWORD,required"`
	DbName         string `env:"DB_NAME,required"`
}

func Init() {
	if err := env.Parse(&Vars); err != nil {
		log.Fatal(err)
	}
}
