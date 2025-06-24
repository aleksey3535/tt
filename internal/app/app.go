package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"task/internal/config"
	"task/internal/handler"
	"task/internal/middleware"
	"task/internal/repository"
	"task/internal/repository/postgres"
	"task/internal/service"
)

type App struct {
	Cfg  *config.Config
	Handler *handler.Handler
}


func New() *App {
	cfg := config.MustLoad()
	log := setupLogger(cfg)
	db := postgres.MustInitDatabase(cfg, log)
	repo := repository.New(db)
	service := service.New(repo, cfg)
	mw := middleware.New(log, service)
	handler := handler.New(log, service, cfg, mw)

	return &App{
		Cfg: 	cfg,
		Handler: handler,
	}
}

func (a *App) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%s", a.Cfg.Port), a.Handler.InitRoutes())
}

func setupLogger(cfg *config.Config) *slog.Logger {
	switch cfg.Env {
	case "local":
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		logger.Info("Logger initialized")
		return logger
	case "prod":
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		logger.Info("Logger initialized")
		return logger
	default:
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		logger.Info("Logger initialized")
		return logger
	}
}