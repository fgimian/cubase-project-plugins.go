package main

import (
	"os"

	"github.com/fgimian/cubase-project-plugins/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
