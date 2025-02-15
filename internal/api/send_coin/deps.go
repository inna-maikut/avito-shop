//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package send_coin

import (
	"context"
)

type coinSending interface {
	Send(ctx context.Context, employeeID int64, targetUsername string, amount int64) error
}
