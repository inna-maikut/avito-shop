package pg

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // register "pgx" database/sql driver
	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-shop/internal/infrastructure/config"
)

func NewDB(ctx context.Context, cfg config.Config) (*sqlx.DB, func(), error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost,
		cfg.DatabasePort, cfg.DatabaseName)

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("sql.Open: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("db.PingContext: %w", err)
	}

	return sqlx.NewDb(db, "pgx"), func() {
		_ = db.Close()
	}, nil
}
