package config

import (
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type appConfig struct {
	Name    string `env:"APP_NAME"`
	Address string `env:"APP_ADDRESS"`
	Debug   bool   `env:"APP_DEBUG"`
}

type dBConfig struct {
	DSN           string `env:"DB_DSN"`
	MaxOpenPool   int    `env:"DB_MAX_OPEN_POOL" envDefault:"10"`
	MaxIdlePool   int    `env:"DB_MAX_IDLE_POOL" envDefault:"10"`
	MaxIdleSecond int    `env:"DB_MAX_IDLE_SECOND" envDefault:"300"`
}

type Config struct {
	App appConfig
	DB  dBConfig
}

func NewConfig(logger *slog.Logger) Config {
	if os.Getenv("APP_ADDRESS") == "" {
		logger.Info("OS Env not found. Try loading .env file")
		if err := godotenv.Load(); err != nil {
			logger.Error("Failed to load .env", slog.Any("error", err))
		}
	}

	cfg := Config{}
	opts := env.Options{RequiredIfNoDef: true}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		logger.Error("Failed to parse config", slog.Any("error", err.Error()))
	}
	return cfg
}
