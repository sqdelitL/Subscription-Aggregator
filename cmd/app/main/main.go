package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sqdelitL/subscription-aggregator/internal/build"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/config"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("config error:", "error", err)
		os.Exit(1)
	}
	log.SlogSetup(cfg.LogLevel)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	app, err := build.NewApp(ctx, cfg)
	if err != nil {
		slog.Error("error with create application", "err", err)
		return
	}
	app.Start()

	<-ctx.Done()
	slog.Debug("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	app.Stop(shutdownCtx)
}
