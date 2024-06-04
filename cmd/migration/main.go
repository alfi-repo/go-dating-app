package main

import (
	"go-dating-app/config"
	migration "go-dating-app/database"
	"go-dating-app/storage"
	"log/slog"
	"os"

	"github.com/pressly/goose/v3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.NewConfig(logger)
	db := storage.NewDB(logger, cfg)

	if err := goose.SetDialect("mysql"); err != nil {
		logger.Error("failed to set dialect", slog.Any("error", err))
	}

	goose.SetBaseFS(migration.EmbedMigrations)

	if err := goose.Up(db, "migration"); err != nil {
		logger.Error("failed to run migrations", slog.Any("error", err))
	}
}
