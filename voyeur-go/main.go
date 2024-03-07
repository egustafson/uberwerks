package main

import (
	"os"

	"github.com/werks/voyeur-go/cmd"
)

// populated by linker at build time
var (
	// GitSummary = git describe --tags --dirty --always
	GitSummary string = "v0.0.0-dirty"

	// BuildDate = date -u +%Y-%m-%dT%H:%M:%SZ
	BuildDate string = "1970-01-01T00:00:00Z"
)

func main() {
	err := cmd.Execute(GitSummary, BuildDate)
	if err != nil {
		// cobra will print an error to stdout/(?)err
		os.Exit(1)
	}
}
