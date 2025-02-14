package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/inna-maikut/avito-shop/internal/api"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-shop/internal/model"
)

type Handler struct {
	authenticating authenticating
}

func New(authenticating authenticating) (*Handler, error) {
	return &Handler{
		authenticating: authenticating,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var authRequest api.AuthRequest
	if ok := api_handler.Parse(r, w, &authRequest); !ok {
		return
	}

	token, err := h.authenticating.Auth(r.Context(), authRequest.Username, authRequest.Password)
	if err != nil {
		if errors.Is(err, model.ErrWrongEmployeePassword) {
			api_handler.Unauthorized(w, "wrong user password")
			return
		}

		fmt.Println("error: ", err)
		api_handler.InternalError(w, "internal server error")
		return
	}

	api_handler.OK(w, api.AuthResponse{
		Token: &token,
	})
}
