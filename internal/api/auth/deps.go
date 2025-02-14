//go:generate mockgen -source deps.go -package $GOPACKAGE -typed -destination mock_deps_test.go
package auth

import (
	"context"
)

type authenticating interface {
	Auth(ctx context.Context, username, password string) (string, error)
}
