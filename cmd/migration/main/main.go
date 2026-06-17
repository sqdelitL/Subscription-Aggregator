package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/config"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/db/migration"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: migration <up|down|new [name]>\n")
		os.Exit(1)
	}
	command := strings.TrimSpace(os.Args[1])
	if command == "new" {
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "usage: migration new <name>\n")
			os.Exit(1)
		}
		name := os.Args[2]
		if err := migration.NewMigrateFile(name); err != nil {
			slog.Error("failed to create migration file", "error", err, "name", name)
			os.Exit(1)
		}
		slog.Info("migration file created", "name", name)
		os.Exit(0)
	}

	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	m, err := migration.New(cfg)
	if err != nil {
		slog.Error("failed to initialize migrations", "error", err)
		os.Exit(1)
	}

	switch command {
	case "up":
		if err := m.Up(); err != nil {
			slog.Error("migration up failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil {
			slog.Error("migration down failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migrations rolled back successfully")
	default:
		slog.Error("unknown command", "command", os.Args[1],
			"expected", "up, down, new")
		os.Exit(1)
	}
}
