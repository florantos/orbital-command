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
	"github.com/go-chi/chi/v5"
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
	auditRepo := repository.NewAuditEventRepo()

	h := handler.NewHandler(logger, moduleRepo, auditRepo)

	r := chi.NewRouter()
	r.Use(middleware.CORS(cfg.CORSAllowedOrigin))

	r.Get("/health", h.Health)
	r.Post("/modules", h.CreateModule)
	r.Get("/modules", h.ReadAllModules)

	r.Get("/audit-events", h.ReadAllAuditEvents)

	logger.Info("Initializing server...", "port", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}

}
