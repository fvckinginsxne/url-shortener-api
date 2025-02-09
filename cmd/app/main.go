package main

import (
	"log"
	"log/slog"
	"os"

	"url-shortener/internal/config"

	"github.com/joho/godotenv"
)

const (
	envPath = "/Users/madw3y/petprojects/url-stortener/.env"

	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("can't load .env file", err)
	}

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO: init logger: slog

	// TODO: init storage: postgresql

	// TODO: init router: chi, chi render

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
