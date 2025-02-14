package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// database
	DatabaseName     string `required:"true" split_words:"true"`
	DatabaseHost     string `required:"true" split_words:"true"`
	DatabasePort     int    `required:"true" split_words:"true"`
	DatabaseUser     string `required:"true" split_words:"true"`
	DatabasePassword string `required:"true" split_words:"true"`

	// http server
	ServerPort int `required:"true" split_words:"true"`
}

func Load() Config {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load(".env")
		if err != nil {
			panic(fmt.Errorf("load godotenv .env config: %w", err))
		}
	}
	if _, err := os.Stat(".env.override"); err == nil {
		err = godotenv.Overload(".env.override")
		if err != nil {
			panic(fmt.Errorf("load godotenv .env.override config: %w", err))
		}
	}

	var cfg Config
	envconfig.MustProcess("", &cfg)
	return cfg
}
