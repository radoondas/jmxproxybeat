package main

import (
	"os"

	"github.com/radoondas/jmxproxybeat/cmd"
)

// Name of this Beat.
var Name = "jmxproxybeat"

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
