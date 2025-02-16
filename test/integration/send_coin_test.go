//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-shop/internal/api"
)

func Test_SendCoin_OK(t *testing.T) {
	setUp()

	username1, username2 := makeUsername(t), makeUsername(t)
	token1, token2 := makeUserToken(t, username1), makeUserToken(t, username2)

	resp := apiPost(t, "/api/sendCoin", token1, api.SendCoinRequest{
		Amount: 200,
		ToUser: username2,
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)

	info1, info2 := getInfo(t, token1), getInfo(t, token2)

	assert.Equal(t, 800, *info1.Coins)
	assert.Equal(t, 1200, *info2.Coins)
}

func Test_SendCoin_JustEnoughBalance(t *testing.T) {
	setUp()

	username1, username2 := makeUsername(t), makeUsername(t)
	token1, token2 := makeUserToken(t, username1), makeUserToken(t, username2)

	resp := apiPost(t, "/api/sendCoin", token1, api.SendCoinRequest{
		Amount: 1000,
		ToUser: username2,
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)

	info1, info2 := getInfo(t, token1), getInfo(t, token2)

	assert.Equal(t, 0, *info1.Coins)
	assert.Equal(t, 2000, *info2.Coins)
}

func Test_SendCoin_NotEnoughBalance(t *testing.T) {
	setUp()

	username1, username2 := makeUsername(t), makeUsername(t)
	token1, token2 := makeUserToken(t, username1), makeUserToken(t, username2)

	resp := apiPost(t, "/api/sendCoin", token1, api.SendCoinRequest{
		Amount: 1200,
		ToUser: username2,
	})
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	info1, info2 := getInfo(t, token1), getInfo(t, token2)

	assert.Equal(t, 1000, *info1.Coins)
	assert.Equal(t, 1000, *info2.Coins)
}
