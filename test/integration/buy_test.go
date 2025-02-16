//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-shop/internal/api"
)

func Test_Buy_OK(t *testing.T) {
	setUp()

	username := makeUsername(t)
	token := makeUserToken(t, username)

	merchName := "powerbank" // price = 200
	resp := apiGet(t, "/api/buy/"+merchName, token)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp = apiGet(t, "/api/buy/"+merchName, token)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	// buy two powerbanks

	info := getInfo(t, token)

	assert.Equal(t, 600, *info.Coins)
	require.Len(t, *info.Inventory, 1)
	assert.Equal(t, *(*info.Inventory)[0].Quantity, 2)
	assert.Equal(t, *(*info.Inventory)[0].Type, merchName)
}

func Test_Buy_NotEnoughBalance(t *testing.T) {
	setUp()

	username := makeUsername(t)
	token := makeUserToken(t, username)

	merchName := "pink-hoody"
	resp := apiGet(t, "/api/buy/"+merchName, token)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp = apiGet(t, "/api/buy/"+merchName, token)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	// buy two pink-hoody, balance is 0 - just enough balance

	resp = apiGet(t, "/api/buy/"+merchName, token)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	// not enough balance

	info := getInfo(t, token)

	assert.Equal(t, 0, *info.Coins)
	require.Len(t, *info.Inventory, 1)
	assert.Equal(t, *(*info.Inventory)[0].Quantity, 2)
	assert.Equal(t, *(*info.Inventory)[0].Type, merchName)
}

func Test_Buy_InvalidMerchName(t *testing.T) {
	setUp()

	username := makeUsername(t)
	token := makeUserToken(t, username)

	merchName := "invalid-merch-name"
	resp := apiGet(t, "/api/buy/"+merchName, token)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	out := parseJSON[api.ErrorResponse](t, resp)

	require.Equal(t, "no merch with name invalid-merch-name", *out.Errors)
}
