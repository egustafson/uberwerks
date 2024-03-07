package voyeurd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/werks/voyeur-go/agent"
)

func Run() error {

	logger := initLogging()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agent := agent.InitAgent(agent.WithLogger(logger))

	// hook signals for shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logger.Infof("received signal: %s", sig.String())
		cancel() // ==> shutdown the agent
	}()

	// start the agent and wait forever
	logger.Info("voyeur agent starting")
	agent.Start(ctx)
	err := agent.AwaitShutdown()
	logger.Info("voyeur agent shutdown")
	if err != nil {
		logger.Warnf("voyeur agent exited with error: %w", err)
	}

	return err
}
