package build

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/sqdelitL/subscription-aggregator/internal/domain/port/repository"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/config"
	database "github.com/sqdelitL/subscription-aggregator/internal/infrastructure/db"
	repoimpl "github.com/sqdelitL/subscription-aggregator/internal/infrastructure/db/repository/subscribe"
	"github.com/sqdelitL/subscription-aggregator/internal/infrastructure/http"
	"github.com/sqdelitL/subscription-aggregator/internal/usecase/subscribe"
)

type App struct {
	db                  *sql.DB
	subscribeWriterRepo repository.SubscriptionWriter
	subscribeReaderRepo repository.SubscriptionReader
	subscribeInteractor *subscribe.Interactor
	server              *http.Server
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	db, err := database.NewConnect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	subscribeWriterRepo := repoimpl.NewSubscriptionWriter(db)
	subscribeReaderRepo := repoimpl.NewSubscriptionReader(db)

	subscribeInteractor := subscribe.New(subscribeWriterRepo, subscribeReaderRepo)

	router := http.NewRouter(subscribeInteractor)
	server := http.NewServer(cfg.ServerPort, router.SetupChi())

	return &App{
		db:                  db,
		subscribeWriterRepo: subscribeWriterRepo,
		subscribeReaderRepo: subscribeReaderRepo,
		subscribeInteractor: subscribeInteractor,
		server:              server,
	}, nil
}

func (app *App) Start() {
	go app.server.Start()
}

func (app *App) Stop(ctx context.Context) {
	if err := app.server.Instance.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
	if err := app.db.Close(); err != nil {
		slog.Error("database close error", "error", err)
	}
}
