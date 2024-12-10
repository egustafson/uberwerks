package voyeurd

import (
	"log/slog"
	"os"
)

func initLogging() *slog.Logger {
	var devFlag = true // evenutally pass in

	logWr := os.Stdout
	level := slog.LevelInfo

	var logger *slog.Logger
	if devFlag {
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

	return logger
}
