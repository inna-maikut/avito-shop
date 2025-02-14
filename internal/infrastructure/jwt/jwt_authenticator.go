package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type tokenProvider interface {
	ParseToken(tokenStr string) (model.TokenInfo, error)
}

type tokenContextKey struct{}

var (
	ErrNoAuthHeader  = errors.New("authorization header is missing")
	ErrClaimsInvalid = errors.New("provided claims do not match expected scopes")
)

// GetJWSFromRequest extracts a JWS string from an Authorization: <jws> header
func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	return authHdr, nil
}

func NewAuthenticator(provider tokenProvider) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(ctx, provider, input)
	}
}

// Authenticate uses the specified validator to ensure a JWT is valid, then makes
// sure that the claims provided by the JWT match the scopes as required in the API.
func Authenticate(ctx context.Context, provider tokenProvider, input *openapi3filter.AuthenticationInput) error {
	// Our security scheme is named BearerAuth, ensure this is the case
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}

	tokenInfo, err := provider.ParseToken(jws)
	if err != nil {
		return fmt.Errorf("validating JWS: %w", err)
	}

	ctx = ContextWithTokenInfo(ctx, tokenInfo)
	input.RequestValidationInput.Request = input.RequestValidationInput.Request.WithContext(ctx)

	return nil
}

func TokenInfoFromContext(ctx context.Context) model.TokenInfo {
	v := ctx.Value(tokenContextKey{})
	res, ok := v.(model.TokenInfo)
	if !ok {
		return model.TokenInfo{}
	}
	return res
}

func ContextWithTokenInfo(ctx context.Context, tokenInfo model.TokenInfo) context.Context {
	return context.WithValue(ctx, tokenContextKey{}, tokenInfo)
}
