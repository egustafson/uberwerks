package server

import (
	"log/slog"
	"os"

	"github.com/egustafson/uberwerks/jsondb-go/server/config"
)

func initServerLogging(config *config.Config) {

	logWr := os.Stdout

	level := slog.LevelInfo
	if config.Flags.Verbose {
		level = slog.LevelDebug
	}

	var logger *slog.Logger
	if config.Flags.DevelMode {
		level = slog.LevelDebug
		logger = slog.New(slog.NewTextHandler(logWr, &slog.HandlerOptions{
			Level: level,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(logWr, &slog.HandlerOptions{
			Level: level,
		}))
	}
	slog.SetDefault(logger)
}
