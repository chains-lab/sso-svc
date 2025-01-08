package main

import (
	"os"

	"github.com/cifra-city/sso-oauth/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
