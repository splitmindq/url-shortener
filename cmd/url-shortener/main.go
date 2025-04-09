package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	_, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to create storage", sl.ErisErr(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	//todo server
}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {

	case envLocal:
		log = setupPrettySlog(slog.LevelDebug)
	case envDev:
		log = setupPrettySlog(slog.LevelDebug)
	case envProd:
		log = setupPrettySlog(slog.LevelInfo)
	}

	return log
}

func setupPrettySlog(s slog.Level) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: s,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
