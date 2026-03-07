package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrations embed.FS

func Migrate(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close() //nolint:errcheck

	goose.SetBaseFS(migrations)
	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("select dialect: %w", err)
	}

	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
