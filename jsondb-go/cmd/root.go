package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "jsondb <sub-command>",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	PersistentPreRunE: initAppHook,
}

var (
	GitSummary string = ""
	BuildDate  string = ""
)

// flags initialized in flags.go

func Execute(gitSummary, buildDate string) error {
	GitSummary = gitSummary
	BuildDate = buildDate
	return rootCmd.Execute()
}

func initAppHook(cmd *cobra.Command, args []string) error {

	//
	// TODO: app initialization
	//

	return nil
}
