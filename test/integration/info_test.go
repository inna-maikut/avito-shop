//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-shop/internal/api"
)

func Test_Info(t *testing.T) {
	setUp()

	username1, username2, username3 := makeUsername(t), makeUsername(t), makeUsername(t)
	token1, token2, token3 := makeUserToken(t, username1), makeUserToken(t, username2), makeUserToken(t, username3)

	// send coin to add a transaction username1 -> username2
	resp := apiPost(t, "/api/sendCoin", token1, api.SendCoinRequest{
		Amount: 300,
		ToUser: username2,
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	// send coin to add a transaction username2 -> username1
	resp = apiPost(t, "/api/sendCoin", token2, api.SendCoinRequest{
		Amount: 100,
		ToUser: username1,
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)
	// send coin to add a transaction username1 -> username3
	resp = apiPost(t, "/api/sendCoin", token3, api.SendCoinRequest{
		Amount: 900,
		ToUser: username1,
	})
	require.Equal(t, http.StatusOK, resp.StatusCode)

	merchName := "pink-hoody"
	resp = apiGet(t, "/api/buy/"+merchName, token1)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	merchName = "pen"
	resp = apiGet(t, "/api/buy/"+merchName, token1)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp = apiGet(t, "/api/buy/"+merchName, token1)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	info := getInfo(t, token1)

	assert.Equal(t, 1180, *info.Coins)

	// check first user coin history
	require.Len(t, *info.CoinHistory.Sent, 1)
	require.Len(t, *info.CoinHistory.Received, 2)
	assert.Equal(t, *(*info.CoinHistory.Sent)[0].ToUser, username2)
	assert.Equal(t, *(*info.CoinHistory.Sent)[0].Amount, 300)
	assert.Equal(t, *(*info.CoinHistory.Received)[0].FromUser, username2)
	assert.Equal(t, *(*info.CoinHistory.Received)[0].Amount, 100)
	assert.Equal(t, *(*info.CoinHistory.Received)[1].FromUser, username3)
	assert.Equal(t, *(*info.CoinHistory.Received)[1].Amount, 900)

	// check first user
	require.Len(t, *info.Inventory, 2)
	assert.Equal(t, *(*info.Inventory)[0].Quantity, 1)
	assert.Equal(t, *(*info.Inventory)[0].Type, "pink-hoody")
	assert.Equal(t, *(*info.Inventory)[1].Quantity, 2)
	assert.Equal(t, *(*info.Inventory)[1].Type, "pen")
}
