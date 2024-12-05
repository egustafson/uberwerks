package cmd

import (
	"log/slog"
	"os"
)

// initLogging initializes logging for cli/client uses of the application.
func initLogging() {

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)
}
