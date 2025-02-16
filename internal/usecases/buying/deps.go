//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package buying

import (
	"context"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type trManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

type employeeRepo interface {
	GetByUsername(ctx context.Context, username string) (*model.Employee, error)
	GetByIDWithLock(ctx context.Context, employeeID int64) (*model.Employee, error)
	IncreaseBalance(ctx context.Context, employeeID, amount int64) error
}

type inventoryRepo interface {
	AddOne(ctx context.Context, employeeID, merchID int64) error
}

type merchRepo interface {
	GetByName(ctx context.Context, name string) (*model.Merch, error)
}
