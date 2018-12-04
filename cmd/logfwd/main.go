package main

import (
	"os"

	"github.com/v3io/logfwd/cmd/logfwd/app"
	"github.com/v3io/go-errors"
)

func main() {
	if err := app.NewCommand().Run(); err != nil {
		errors.PrintErrorStack(os.Stderr, err, 5)
		os.Exit(1)
	}

	os.Exit(0)
}
