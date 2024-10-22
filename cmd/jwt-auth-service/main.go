package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"jwt-auth-service/internal/config"
	"jwt-auth-service/internal/http_server/handlers/url/auth"
	"jwt-auth-service/internal/http_server/handlers/url/registration"
	"jwt-auth-service/internal/http_server/middleware/logger"
	"jwt-auth-service/internal/storage/postgresql"
	"log/slog"
	"net/http"
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
	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.Logger,
		logger.New(log),
		middleware.Recoverer,
		middleware.URLFormat,
	)

	router.Post("/register", registration.New(log, storage))
	router.Get("/auth", auth.New(log, storage))

	// TODO: RUN SERVER
	log.Info("starting server...", slog.String("address", cfg.HTTPServer.Host+":"+cfg.HTTPServer.Port))
	serverAddr := cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
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
