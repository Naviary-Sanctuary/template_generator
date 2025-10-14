package main

import (
	"os"

	"github.com/Naviary-Sanctuary/template_generator/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
