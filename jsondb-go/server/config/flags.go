package config

import "github.com/spf13/pflag"

func parseFlags(args []string, flags *Flags) error {
	fs := pflag.FlagSet{}
	fs.BoolVar(&(flags.DevelMode), "dev", false, "development mode")

	return fs.Parse(args)
}
