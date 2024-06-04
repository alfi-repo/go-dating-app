package storage

import (
	"database/sql"
	"go-dating-app/config"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(logger *slog.Logger, config config.Config) *sql.DB {
	db, err := sql.Open("mysql", config.DB.DSN)
	if err != nil {
		logger.Error("Error when opening DB", slog.Any("error", err))
		return nil
	}

	db.SetMaxOpenConns(config.DB.MaxOpenPool)
	db.SetMaxIdleConns(config.DB.MaxIdlePool)
	db.SetConnMaxLifetime(time.Duration(config.DB.MaxIdleSecond) * time.Second)
	return db
}
