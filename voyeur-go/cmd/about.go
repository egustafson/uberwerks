package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "about voyeur executable",
	RunE:  doAbout,
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}

func doAbout(cmd *cobra.Command, args []string) error {

	// TODO:  investigate the CLI "output" method
	//

	// output is valid YAML
	fmt.Println("---")
	fmt.Printf("voyeur-version: %s\n", GitSummary)
	fmt.Printf("build-date: %s\n", BuildDate)
	fmt.Println("...")

	return nil
}
