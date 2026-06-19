package main

import (
	"context"
	"errors"
	"os"

	"github.com/satoricorp/gx/internal/cli"
)

func main() {
	if err := cli.Execute(context.Background()); err != nil && !errors.Is(err, context.Canceled) {
		os.Exit(1)
	}
}

func shouldLaunch(args []string) bool {
	return cli.ShouldLaunch(args)
}
