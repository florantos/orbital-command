package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	LogLevel    string
	DatabaseURL string
}

func Load() (Config, error) {
	_ = godotenv.Load() // ignore this error: no env file in prod

	port, ok := os.LookupEnv("PORT")
	if !ok {
		return Config{}, fmt.Errorf("cannot find env variable PORT")
	}
	if port == "" {
		return Config{}, fmt.Errorf("env variable PORT is empty")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return Config{}, fmt.Errorf("cannot find env variable DATABASE_URL")
	}
	if url == "" {
		return Config{}, fmt.Errorf("env variable DATABASE_URL is empty")
	}

	cfg := Config{
		Port:        port,
		LogLevel:    logLevel,
		DatabaseURL: url,
	}

	return cfg, nil
}
