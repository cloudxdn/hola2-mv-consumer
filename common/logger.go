package common

import (
	"log"
	"log/slog"
	"os"
)

func SetupLogger(logLevel string) {
	var level slog.Level

	switch logLevel {
	case "debug":
		level = slog.LevelDebug
		log.Println("Staring Debug mode")
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
}
