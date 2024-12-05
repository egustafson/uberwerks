// Package jsondb implements the JSON-DB daemon.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/egustafson/uberwerks/jsondb-go/server/api"
	"github.com/egustafson/uberwerks/jsondb-go/server/config"
)

func Start(args []string, flags config.Flags) (err error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// hook signals for shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		slog.Info(fmt.Sprintf("received signal: %s", sig.String()))
		cancel()
	}()

	var cfg *config.Config
	if cfg, ctx, err = config.InitConfig(ctx, args, flags); err != nil {
		return
	}
	initServerLogging(cfg)

	// TODO: remainder of daemon and api start-up

	api.Run(ctx, cfg)
	<-ctx.Done() // block until the context is canceled
	return nil
}
