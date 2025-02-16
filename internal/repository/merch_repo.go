package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type MerchRepository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewMerchRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) (*MerchRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if getter == nil {
		return nil, errors.New("getter is nil")
	}

	return &MerchRepository{
		db:     db,
		getter: getter,
	}, nil
}

func (r *MerchRepository) trOrDB(ctx context.Context) trmsqlx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db)
}

func (r *MerchRepository) GetByName(ctx context.Context, name string) (*model.Merch, error) {
	var merch Merch

	q := "SELECT id, name, price FROM merch WHERE name = $1"

	err := r.trOrDB(ctx).GetContext(ctx, &merch, q, name)
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
