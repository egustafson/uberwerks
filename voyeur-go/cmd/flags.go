package cmd

// flags for voyeur

var (
	verboseFlag bool = false
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false,
		"verbose output")
}
