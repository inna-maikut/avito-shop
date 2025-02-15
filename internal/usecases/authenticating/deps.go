//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package authenticating

import (
	"context"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type employeeRepo interface {
	GetByUsername(ctx context.Context, username string) (*model.Employee, error)
	Create(ctx context.Context, username, passwordHash string, balance int64) (*model.Employee, error)
}

type tokenProvider interface {
	CreateToken(username string, userID int64) (string, error)
}
