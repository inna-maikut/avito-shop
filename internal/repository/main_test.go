//go:build integration

package repository

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib" // register "pgx" database/sql driver
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-shop/internal/infrastructure/config"
)

func setUp(t *testing.T) *sqlx.DB {
	cfg := config.Load()

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost,
		cfg.DatabasePort, cfg.DatabaseName)

	db, err := sql.Open("pgx", databaseURL)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = db.Close()
	})

	return sqlx.NewDb(db, "pgx")
}
