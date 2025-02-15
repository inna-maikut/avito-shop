package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type MerchRepository struct {
	db *sqlx.DB
}

func NewMerchRepository(db *sqlx.DB) (*MerchRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &MerchRepository{
		db: db,
	}, nil
}

func (r *MerchRepository) GetByName(ctx context.Context, name string) (*model.Merch, error) {
	var merch Merch

	err := r.db.GetContext(ctx, &merch, "SELECT * FROM merch WHERE name = $1", name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrMerchNotFound
		}
		return nil, fmt.Errorf("db.GetContext: %w", err)
	}

	return &model.Merch{
		ID:    merch.ID,
		Name:  merch.Name,
		Price: merch.Price,
	}, nil
}
