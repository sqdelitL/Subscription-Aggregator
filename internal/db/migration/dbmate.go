package migration

import (
	"embed"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/amacneil/dbmate/v2/pkg/dbutil"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/sqdelitL/subscription-aggregator/internal/config"
)

const migrationsDir = "internal/db/migration/sql"

//go:embed sql
var migrationsFS embed.FS

type Migrator struct {
	db *dbmate.DB
}

func New(cfg *config.Config) (*Migrator, error) {
	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.DataBase.UserName, cfg.DataBase.Password),
		Host:     fmt.Sprintf("%s:%s", cfg.DataBase.Host, cfg.DataBase.Port),
		Path:     cfg.DataBase.Name,
		RawQuery: "sslmode=disable",
	}
	db := dbmate.New(u)
	db.Strict = true // при true упадет миграция если в неправильном порядке применена
	db.FS = migrationsFS
	db.MigrationsDir = []string{"sql"}
	return &Migrator{
		db: db,
	}, nil
}

func (m *Migrator) Up() error {
	err := m.db.Migrate()
	if err != nil {
		return err
	}
	return nil
}

func (m *Migrator) Down() error {
	err := m.db.Rollback()
	if err != nil {
		return err
	}
	return nil
}

const migrationTemplate = "-- migrate:up\n\n\n-- migrate:down\n\n"

func NewMigrateFile(fileName string) error {
	timestamp := time.Now().UTC().Format("20060102150405")
	if fileName == "" {
		return dbmate.ErrNoMigrationName
	}
	fileName = fmt.Sprintf("%s_%s.sql", timestamp, fileName)

	if err := ensureDir(migrationsDir); err != nil {
		return err
	}

	path := filepath.Join(migrationsDir, fileName)
	slog.Info("Creating migration", "path", path)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return dbmate.ErrMigrationAlreadyExist
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer dbutil.MustClose(file)
	_, err = file.WriteString(migrationTemplate)
	return err
}

func ensureDir(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("%w `%s`", dbmate.ErrCreateDirectory, dir)
	}

	return nil
}
