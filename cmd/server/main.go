package main

import (
	"log"

	"github.com/florantos/orbital-command/internal/config"
	"github.com/florantos/orbital-command/internal/logger"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log := logger.New(cfg.LogLevel, cfg.Env)
	log.Info("Initializing server...")
	log.Debug("temp debug")
}
