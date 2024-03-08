package cmd

import (
	"github.com/spf13/cobra"

	"github.com/egustafson/uberwerks/jsondb-go/server"
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

	err := server.Start()
	return err
}
