package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/lib/pq"
	"github.com/sqdelitL/subscription-aggregator/internal/config"
)

const driver = "postgres"

func NewConnect(ctx context.Context, cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn(cfg))
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.DataBase.MaxIdleConns)
	db.SetMaxOpenConns(cfg.DataBase.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(cfg.DataBase.ConnMaxLifetime))

	go pingWithRetry(ctx, db, cfg)
	return db, nil
}

func pingWithRetry(ctx context.Context, db *sql.DB, cfg *config.Config) {
	retry := cfg.DataBase.Retry

	retryInterval := time.Duration(retry.RetryDelaySeconds) * time.Second
	maxInterval := time.Duration(retry.MaxIntervalSeconds) * time.Second
	maxRetries := retry.MaxRetries
	backoffMultiplier := retry.BackoffMultiplierSeconds

	operation := func() error {
		if pingErr := db.Ping(); pingErr != nil {
			slog.Warn("failed to ping the database", "error", pingErr)
			return pingErr
		}

		slog.Info("successful database ping")
		return nil
	}

	var b backoff.BackOff

	if backoffMultiplier > 0 {
		exp := backoff.NewExponentialBackOff()
		exp.InitialInterval = retryInterval
		exp.Multiplier = backoffMultiplier
		exp.MaxInterval = maxInterval
		exp.MaxElapsedTime = 0

		if maxRetries > 0 {
			b = backoff.WithMaxRetries(exp, uint64(maxRetries))
		} else {
			b = exp
		}
	} else {
		fb := backoff.NewConstantBackOff(retryInterval)
		if maxRetries > 0 {
			b = backoff.WithMaxRetries(fb, uint64(maxRetries))
		} else {
			b = fb
		}
	}
	b = backoff.WithContext(b, ctx)

	if err := backoff.Retry(operation, b); err != nil {
		slog.Error("failed to connect to the database after all attempts.", "error", err)
	}
}

func dsn(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DataBase.Host, cfg.DataBase.Port, cfg.DataBase.UserName, cfg.DataBase.Password, cfg.DataBase.Name)
}
