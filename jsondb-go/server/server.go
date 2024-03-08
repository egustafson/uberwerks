// Package jsondb implements the JSON-DB daemon.
package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Masterminds/log-go"

	"github.com/egustafson/uberwerks/jsondb-go/server/api"
)

func Start() error {

	//initLogging()  // TODO:  belongs in another package, ?? mx possibly?

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// hook signals for shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Infof("received signal: %s", sig.String())
		cancel()
	}()

	// TODO: remainder of daemon and api start-up

	api.Run(ctx)
	<-ctx.Done() // block until the context is canceled
	return nil
}
