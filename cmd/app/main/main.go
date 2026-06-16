package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/sqdelitL/subscription-aggregator/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("config:", "error", err)
		os.Exit(1)
	}
	fmt.Println(cfg)

}
