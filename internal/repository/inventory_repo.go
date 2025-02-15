package repository

import (
	"context"
	"errors"
	"fmt"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type InventoryRepository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewInventoryRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) (*InventoryRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if getter == nil {
		return nil, errors.New("getter is nil")
	}

	return &InventoryRepository{
		db:     db,
		getter: getter,
	}, nil
}

func (r *InventoryRepository) trOrDB(ctx context.Context) trmsqlx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db)
}

func (r *InventoryRepository) GetByEmployee(ctx context.Context, employeeID int64) ([]model.Inventory, error) {
	var inventories []InventoryWithMerchName

	q := "SELECT i.employee_id, i.merch_id, i.quantity, merch.name as merch_name " +
		"FROM inventory i " +
		"INNER JOIN merch on merch.id = i.merch_id " +
		"WHERE employee_id = $1"

	err := r.trOrDB(ctx).SelectContext(ctx, &inventories, q, employeeID)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	res := make([]model.Inventory, 0, len(inventories))
	for _, inventory := range inventories {
		res = append(res, model.Inventory{
			EmployeeID: inventory.EmployeeID,
			MerchID:    inventory.MerchID,
			Quantity:   inventory.Quantity,
			MerchName:  inventory.MerchName,
		})
	}

	return res, nil
}
