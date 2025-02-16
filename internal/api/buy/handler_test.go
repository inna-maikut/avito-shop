package buy

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
	buyingMock := NewMockbuying(ctrl)

	buyingMock.EXPECT().
		Buy(gomock.Any(), int64(1234), "socks").
		Return(nil)

	handler, err := New(buyingMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/buy/socks", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("merchName", "socks")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		EmployeeID: 1234,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Handle_ErrMerchNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	buyingMock := NewMockbuying(ctrl)

	buyingMock.EXPECT().
		Buy(gomock.Any(), int64(1234), "no-such-merch").
		Return(model.ErrMerchNotFound)

	handler, err := New(buyingMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/buy/no-such-merch", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("merchName", "no-such-merch")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		EmployeeID: 1234,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var response api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "no merch with name no-such-merch", *response.Errors)
}

func TestHandler_Handle_ErrNotEnoughBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	buyingMock := NewMockbuying(ctrl)

	buyingMock.EXPECT().
		Buy(gomock.Any(), int64(1234), "socks").
		Return(model.ErrNotEnoughBalance)

	handler, err := New(buyingMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/buy/socks", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("merchName", "socks")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		EmployeeID: 1234,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var response api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "not enough balance", *response.Errors)
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	buyingMock := NewMockbuying(ctrl)

	buyingMock.EXPECT().
		Buy(gomock.Any(), int64(1234), "socks").
		Return(assert.AnError)

	handler, err := New(buyingMock, zap.NewNop())
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/buy/socks", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("merchName", "socks")
	req = req.WithContext(jwt.ContextWithTokenInfo(req.Context(), model.TokenInfo{
		EmployeeID: 1234,
	}))
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	var response api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "internal server error", *response.Errors)
}
