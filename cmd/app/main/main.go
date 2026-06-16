package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/sqdelitL/subscription-aggregator/internal/config"
	"github.com/sqdelitL/subscription-aggregator/internal/db/migration"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "run database migrations before start")
	flag.Parse()

	cfg, err := config.New()
	if err != nil {
		slog.Error("config:", "error", err)
		os.Exit(1)
	}

	if *migrateFlag {
		m, err := migration.New(cfg)
		if err != nil {
			slog.Error("failed to create migrator", "error", err)
			os.Exit(1)
		}
		if err := m.Up(); err != nil {
			slog.Error("failed to run migrations", "error", err)
			os.Exit(1)
		}
		slog.Info("migrations applied successfully")
	}
}
