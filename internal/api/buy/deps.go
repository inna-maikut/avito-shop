//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package buy

import (
	"context"
)

type buying interface {
	Buy(ctx context.Context, employeeID int64, merchName string) error
}
