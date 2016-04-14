package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/radoondas/jmxbeat/beater"
)

func main() {
	err := beat.Run("jmxbeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
