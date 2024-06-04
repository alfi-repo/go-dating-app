package main

import (
	"go-dating-app/api/rest"
	"go-dating-app/app/repository"
	"go-dating-app/app/service"
	"go-dating-app/config"
	"go-dating-app/storage"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Bootstrap.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.NewConfig(logger)
	db := storage.NewDB(logger, cfg)

	// Repository.
	userRepo := repository.NewUserRepository(db)

	// Service.
	authService := service.NewAuthService(&userRepo)

	// API Handler.
	restHandler := rest.NewHandler(logger, authService)

	// API Rest.
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routeAuth := e.Group("/auth")
	routeAuth.POST("/login", restHandler.Login)
	routeAuth.POST("/register", restHandler.Registration)

	e.Logger.Fatal(e.Start(cfg.App.Address))
}
