//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/inna-maikut/avito-shop/internal/api"
)

func Test_AuthMW_NoAuth(t *testing.T) {
	setUp()

	wantPlainError := "security requirements failed: getting jws: authorization header is missing\n"

	t.Run("GET /api/buy/{merchName}", func(t *testing.T) {
		resp := apiGet(t, "/api/buy/pen", "")
		assertResponsePlainError(t, resp, http.StatusUnauthorized, wantPlainError)
	})

	t.Run("POST /api/sendCoin", func(t *testing.T) {
		resp := apiPost(t, "/api/sendCoin", "", api.SendCoinRequest{})
		assertResponsePlainError(t, resp, http.StatusUnauthorized, wantPlainError)
	})

	t.Run("GET /api/info", func(t *testing.T) {
		resp := apiGet(t, "/api/info", "")
		assertResponsePlainError(t, resp, http.StatusUnauthorized, wantPlainError)
	})
}

func Test_AuthMW_InvalidToken(t *testing.T) {
	setUp()

	wantPlainError := "security requirements failed: validating JWS: jwt.Parse: token contains an invalid number of segments\n"

	t.Run("GET /api/buy/{merchName}", func(t *testing.T) {
		resp := apiGet(t, "/api/buy/pen", "1234")
		assertResponsePlainError(t, resp, http.StatusUnauthorized, wantPlainError)
	})

	t.Run("POST /api/sendCoin", func(t *testing.T) {
		resp := apiPost(t, "/api/sendCoin", "1234", api.SendCoinRequest{})
		assertResponsePlainError(t, resp, http.StatusUnauthorized, wantPlainError)
	})

	t.Run("GET /api/info", func(t *testing.T) {
		resp := apiGet(t, "/api/info", "1234")
		assertResponsePlainError(t, resp, http.StatusUnauthorized, wantPlainError)
	})
}

func Test_AuthMW_InvalidSignature(t *testing.T) {
	setUp()

	token := makeUserToken(t, makeUsername(t))
	invalidToken := token + "1"

	wantPlainError := "security requirements failed: validating JWS: jwt.Parse: signature is invalid\n"

	t.Run("GET /api/buy/{merchName}", func(t *testing.T) {
		resp := apiGet(t, "/api/buy/pen", invalidToken)
		assertResponsePlainError(t, resp, http.StatusUnauthorized, wantPlainError)
	})

	t.Run("POST /api/sendCoin", func(t *testing.T) {
		resp := apiPost(t, "/api/sendCoin", invalidToken, api.SendCoinRequest{})
		assertResponsePlainError(t, resp, http.StatusUnauthorized, wantPlainError)
	})

	t.Run("GET /api/info", func(t *testing.T) {
		resp := apiGet(t, "/api/info", invalidToken)
		assertResponsePlainError(t, resp, http.StatusUnauthorized, wantPlainError)
	})
}
