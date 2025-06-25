package main

import (
	"os"

	"github.com/chains-lab/chains-auth/cmd/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
