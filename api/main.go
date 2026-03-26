package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"zhacked.me/oxyl/api/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := cmd.Execute(ctx); err != nil {
		panic(err)
	}
}
