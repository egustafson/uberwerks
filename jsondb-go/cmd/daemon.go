package cmd

import (
	"github.com/spf13/cobra"

	"github.com/egustafson/uberwerks/jsondb-go/server"
	"github.com/egustafson/uberwerks/jsondb-go/server/config"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "start as a daemon",
	RunE:  doDaemon,
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}

func doDaemon(cmd *cobra.Command, args []string) error {

	config := &config.Config{
		//
		// TODO: dynamically load/populate config
		//
		DSN:  ":memory:",
		Port: 8080,
	}

	err := server.Start(config)
	return err
}
