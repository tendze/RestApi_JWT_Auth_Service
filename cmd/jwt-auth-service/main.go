package main

import (
	"fmt"
	"jwt-auth-service/internal/config"
	"jwt-auth-service/internal/storage/postgresql"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: LOAD CONFIG
	cfg := config.MustLoad()

	// TODO: INIT LOGGER
	log := setupLogger(cfg.Env)

	// TODO: INIT STORAGE: POSTGRESQL
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.Username,
		cfg.DB.DBPassword,
		cfg.DB.Host, cfg.DB.Port,
		cfg.DB.DBName,
		cfg.DB.SSLMode,
	)
	storage, err := postgresql.New(dsn)
	if err != nil {
		log.Error("failed to init storage:", err)
	}
	defer storage.DB.Close()

	// TODO: INIT ROUTER CHI:"chi render"

	// TODO: RUN SERVER
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
