//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package coin_sending

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

type transactionRepo interface {
	Add(ctx context.Context, senderID, receiverID, amount int64) error
}
