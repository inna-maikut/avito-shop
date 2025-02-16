//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-shop/internal/api"
)

func Test_Auth_RegisterAndLogin(t *testing.T) {
	setUp()

	username := makeUsername(t)
	// register
	token1 := makeUserToken(t, username)

	// login
	token2 := makeUserToken(t, username)

	// token1 can be either equal to token2 or not, depends on timings ("exp" can be different in JWT claims)

	// buy a pen with token2
	merchName := "pen" // price = 10
	resp := apiGet(t, "/api/buy/"+merchName, token2)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// check changes with first token1
	info := getInfo(t, token1)

	// identify that employee is same in token1 and token2
	assert.Equal(t, 990, *info.Coins)
	require.Len(t, *info.Inventory, 1)
	assert.Equal(t, *(*info.Inventory)[0].Quantity, 1)
	assert.Equal(t, *(*info.Inventory)[0].Type, merchName)
}

func Test_Auth_LoginWrongPassword(t *testing.T) {
	setUp()

	username := makeUsername(t)
	// register
	_ = makeUserToken(t, username)

	// login with wrong password
	resp := apiPost(t, "/api/auth", "", api.AuthRequest{
		Username: username,
		Password: password + "-wrong",
	})

	assertResponseError(t, resp, http.StatusUnauthorized, "wrong user password")
}
