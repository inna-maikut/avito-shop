package buy

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/inna-maikut/avito-shop/internal"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-shop/internal/model"
)

type Handler struct {
	buying buying
	logger internal.Logger
}

func New(buying buying, logger internal.Logger) (*Handler, error) {
	if buying == nil {
		return nil, errors.New("buying is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		buying: buying,
		logger: logger,
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

		err = fmt.Errorf("buying.Buy: %w", err)
		h.logger.Error("GET /api/buy/{merchName} internal error", zap.Error(err),
			zap.String("merchName", merchName))
		api_handler.InternalError(w, "internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}
