package cmd

import (
	"github.com/egustafson/uberwerks/jsondb-go/jsondbd"
	"github.com/spf13/cobra"
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

	err := jsondbd.Run()
	return err
}
