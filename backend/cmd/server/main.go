package main

import (
	"log"
	"net/http"
	"os"

	"github.com/florantos/orbital-command/internal/config"
	"github.com/florantos/orbital-command/internal/handler"
	applogger "github.com/florantos/orbital-command/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := applogger.New(cfg.LogLevel, cfg.Env)
	h := handler.NewHandler(logger, nil)

	http.HandleFunc("/health", h.Health)

	logger.Info("Initializing server...", "port", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
