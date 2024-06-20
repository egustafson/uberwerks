package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var healthzCmd = &cobra.Command{
	Use:   "healthz",
	Short: "check health of the server",
	RunE:  doHealthz,
}

func init() {
	rootCmd.AddCommand(healthzCmd)
}

func doHealthz(cmd *cobra.Command, args []string) error {

	// TODO:  investigate the CLI "output" method

	fmt.Println("health -- stub")
	return nil
}
