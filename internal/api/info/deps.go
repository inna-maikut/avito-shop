//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package info

import (
	"context"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type infoCollecting interface {
	Collect(ctx context.Context, employeeID int64) (model.EmployeeInfo, error)
}
