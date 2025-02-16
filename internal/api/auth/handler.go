package auth

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/inna-maikut/avito-shop/internal"
	"github.com/inna-maikut/avito-shop/internal/api"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-shop/internal/model"
)

type Handler struct {
	authenticating authenticating
	logger         internal.Logger
}

func New(authenticating authenticating, logger internal.Logger) (*Handler, error) {
	if authenticating == nil {
		return nil, errors.New("authenticating is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		authenticating: authenticating,
		logger:         logger,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var authRequest api.AuthRequest
	if ok := api_handler.Parse(r, w, &authRequest); !ok {
		return
	}

	if authRequest.Username == "" || len(authRequest.Username) > 1024 {
		api_handler.BadRequest(w, "username should contain at least one character and no more than 1024 bytes")
		return
	}

	if authRequest.Password == "" || len(authRequest.Password) > 1024 {
		api_handler.BadRequest(w, "password should contain at least one character and no more than 1024 bytes")
		return
	}

	token, err := h.authenticating.Auth(r.Context(), authRequest.Username, authRequest.Password)
	if err != nil {
		if errors.Is(err, model.ErrWrongEmployeePassword) {
			api_handler.Unauthorized(w, "wrong user password")
			return
		}

		err = fmt.Errorf("authenticating.Auth: %w", err)
		h.logger.Error("POST /api/auth internal error", zap.Error(err), zap.Any("request", authRequest))
		api_handler.InternalError(w, "internal server error")
		return
	}

	api_handler.OK(w, api.AuthResponse{
		Token: &token,
	})
}
