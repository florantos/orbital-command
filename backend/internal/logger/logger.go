package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func New(logLevel, env string) *slog.Logger {
	var level slog.Level

	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		fmt.Printf("%s is not a valid log level, defaulting to 'info'\n", logLevel)
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
