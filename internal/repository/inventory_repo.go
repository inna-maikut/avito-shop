package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) (*InventoryRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &InventoryRepository{
		db: db,
	}, nil
}

func (r *InventoryRepository) GetByEmployee(ctx context.Context, employeeID int64) ([]model.Inventory, error) {
	var inventories []InventoryWithMerchName

	q := "SELECT i.employee_id, i.merch_id, i.quantity, merch.name as merch_name " +
		"FROM inventory i " +
		"INNER JOIN merch on merch.id = i.merch_id " +
		"WHERE employee_id = $1"

	err := r.db.SelectContext(ctx, &inventories, q, employeeID)
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
