package cmd

// global flags for jsondb

var (
	verboseFlag bool = false
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false,
		"verbose output")
}
