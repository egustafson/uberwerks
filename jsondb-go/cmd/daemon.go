package cmd

import (
	"github.com/spf13/cobra"

	"github.com/egustafson/uberwerks/jsondb-go/server"
	"github.com/egustafson/uberwerks/jsondb-go/server/config"
)

var daemonCmd = &cobra.Command{
	Use:                "daemon",
	Short:              "start as a daemon",
	DisableFlagParsing: true, // flags parsed in server
	RunE:               doDaemon,
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}

func doDaemon(cmd *cobra.Command, args []string) error {
	flags := config.Flags{
		Verbose: verboseFlag,
	}
	return server.Start(args, flags)
}
