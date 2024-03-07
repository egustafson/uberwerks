package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "voyeur <sub-command>",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Short:             "uber executable for voyeur",
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

func initAppHook(cmd *cobra.Command, args []string) (err error) {

	//
	// app initialization
	//

	return nil
}
