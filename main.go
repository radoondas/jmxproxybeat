package main

import (
	"os"

	"github.com/radoondas/jmxproxybeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
