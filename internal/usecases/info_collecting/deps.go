//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package info_collecting

import (
	"context"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type employeeRepo interface {
	GetByID(ctx context.Context, employeeID int64) (*model.Employee, error)
}

type transactionRepo interface {
	GetByEmployee(ctx context.Context, employeeID int64) ([]model.Transaction, error)
}

type inventoryRepo interface {
	GetByEmployee(ctx context.Context, employeeID int64) ([]model.Inventory, error)
}
