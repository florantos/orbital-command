package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/florantos/orbital-command/internal/config"
	"github.com/florantos/orbital-command/internal/handler"
	applogger "github.com/florantos/orbital-command/internal/logger"
	"github.com/florantos/orbital-command/internal/middleware"
	"github.com/florantos/orbital-command/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := applogger.New(cfg.LogLevel, cfg.Env)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	moduleRepo := repository.NewModuleRepo(pool)
	auditRepo := repository.NewAuditEventRepo(pool)

	h := handler.NewHandler(logger, moduleRepo, auditRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/modules", h.CreateModule)

	logger.Info("Initializing server...", "port", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, middleware.CORS(cfg.CORSAllowedOrigin)(mux)); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
