package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/inna-maikut/avito-shop/internal/api"
	"github.com/inna-maikut/avito-shop/internal/model"
)

func TestHandler_Handle_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	authenticatingMock.EXPECT().
		Auth(gomock.Any(), "test1", "password1").
		Return("token1", nil)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"username": "test1", "password": "password1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var response api.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "token1", *response.Token)
}

func TestHandler_Handle_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	authenticatingMock.EXPECT().
		Auth(gomock.Any(), "test1", "password1").
		Return("", model.ErrWrongEmployeePassword)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"username": "test1", "password": "password1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	var response api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "wrong user password", *response.Errors)
}

func TestHandler_Handle_UsernameEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"username": "", "password": "password1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var response api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "username should contain at least one character and no more than 1024 bytes", *response.Errors)
}

func TestHandler_Handle_UsernameBigLen(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"username": "` + strings.Repeat("a", 1025) + `", "password": "password1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var response api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "username should contain at least one character and no more than 1024 bytes", *response.Errors)
}

func TestHandler_Handle_PasswordEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"username": "test1", "password": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var response api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "password should contain at least one character and no more than 1024 bytes", *response.Errors)
}

func TestHandler_Handle_PasswordBigLen(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"username": "test1", "password": "` + strings.Repeat("a", 1025) + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var response api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "password should contain at least one character and no more than 1024 bytes", *response.Errors)
}

func TestHandler_Handle_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	authenticatingMock := NewMockauthenticating(ctrl)

	authenticatingMock.EXPECT().
		Auth(gomock.Any(), "test1", "password1").
		Return("", assert.AnError)

	handler, err := New(authenticatingMock, zap.NewNop())
	require.NoError(t, err)

	validData := []byte(`{"username": "test1", "password": "password1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	var response api.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, "internal server error", *response.Errors)
}
