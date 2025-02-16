package info

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-shop/internal/api"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-shop/internal/model"
)

func TestHandler_Handle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	infoCollectingMock := NewMockinfoCollecting(ctrl)

	infoCollectingMock.EXPECT().
		Collect(gomock.Any(), int64(1001)).
		Return(model.EmployeeInfo{
			Coins: 100500,
			Inventory: []model.Inventory{
				{
					EmployeeID: 1001,
					MerchID:    1,
					Quantity:   10,
					MerchName:  "socks",
				},
			},
			ReceivedTransactions: []model.Transaction{
				{
					IsSender:               false,
					CounterpartyEmployeeID: 1005,
					CounterpartyUsername:   "test2",
					Amount:                 300,
				},
			},
			SentTransactions: []model.Transaction{
				{
					IsSender:               true,
					CounterpartyEmployeeID: 1006,
					CounterpartyUsername:   "test3",
					Amount:                 500,
				},
			},
		}, nil)

	handler, err := New(infoCollectingMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/info", nil)
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		EmployeeID: 1001,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `
	{
		"coins": 100500,
		"inventory": [
			{
				"quantity": 10,
				"type": "socks"
			}
		],
		"coinHistory": {
			"received": [
				{
					"fromUser": "test2",
					"amount": 300
				}
			],
			"sent": [
				{
					"toUser": "test3",
					"amount": 500
				}
			]
		}
	}`, w.Body.String())
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	infoCollectingMock := NewMockinfoCollecting(ctrl)

	infoCollectingMock.EXPECT().
		Collect(gomock.Any(), int64(1001)).
		Return(model.EmployeeInfo{}, assert.AnError)

	handler, err := New(infoCollectingMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		EmployeeID: 1001,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	var response api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "internal server error", *response.Errors)
}
