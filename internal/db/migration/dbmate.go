package migration

import (
	"embed"
	"fmt"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/sqdelitL/subscription-aggregator/internal/config"
)

const migrationsDir = "sql"

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
	db.MigrationsDir = []string{migrationsDir}
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

func (m *Migrator) NewMigrateFile(fileName string) error {
	return m.db.NewMigration(fileName)
}
