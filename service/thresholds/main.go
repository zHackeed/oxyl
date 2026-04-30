package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"zhacked.me/oxyl/service/thresholds/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := cmd.Execute(ctx); err != nil {
		slog.Error("failed to start up", "error", err)
		os.Exit(1)
	}
}
