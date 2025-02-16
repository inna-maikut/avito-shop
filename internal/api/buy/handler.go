package buy

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/inna-maikut/avito-shop/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-shop/internal/model"
)

type Handler struct {
	buying buying
}

func New(buying buying) (*Handler, error) {
	return &Handler{
		buying: buying,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenInfo := jwt.TokenInfoFromContext(r.Context())

	merchName := r.PathValue("merchName")
	if merchName == "" {
		api_handler.BadRequest(w, "merchName is required")
		return
	}

	err := h.buying.Buy(ctx, tokenInfo.EmployeeID, merchName)
	if err != nil {
		if errors.Is(err, model.ErrMerchNotFound) {
			api_handler.BadRequest(w, "no merch with name "+merchName)
			return
		}
		if errors.Is(err, model.ErrNotEnoughBalance) {
			api_handler.BadRequest(w, "not enough balance")
			return
		}

		fmt.Println(err)
		api_handler.InternalError(w, "internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}
