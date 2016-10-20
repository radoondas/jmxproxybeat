package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/radoondas/jmxproxybeat/beater"
)

// Name of this Beat.
var Name = "jmxproxybeat"

func main() {
	err := beat.Run(Name, "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
