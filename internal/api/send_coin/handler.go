package send_coin

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/inna-maikut/avito-shop/internal"
	"github.com/inna-maikut/avito-shop/internal/api"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-shop/internal/model"
)

type Handler struct {
	coinSending coinSending
	logger      internal.Logger
}

func New(coinSending coinSending, logger internal.Logger) (*Handler, error) {
	if coinSending == nil {
		return nil, errors.New("coinSending is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}
	return &Handler{
		coinSending: coinSending,
		logger:      logger,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenInfo := jwt.TokenInfoFromContext(r.Context())

	var sendCoinRequest api.SendCoinRequest
	if ok := api_handler.Parse(r, w, &sendCoinRequest); !ok {
		return
	}

	err := h.coinSending.Send(ctx, tokenInfo.EmployeeID, sendCoinRequest.ToUser, int64(sendCoinRequest.Amount))
	if err != nil {
		if errors.Is(err, model.ErrSendingCoinsToMyselfNotAllowed) {
			api_handler.BadRequest(w, "sending coins to yourself not allowed")
			return
		}
		if errors.Is(err, model.ErrNotEnoughBalance) {
			api_handler.BadRequest(w, "not enough balance")
			return
		}

		err = fmt.Errorf("coinSending.Send: %w", err)
		h.logger.Error("GET /api/sendCoin internal error", zap.Error(err), zap.Any("tokenInfo", tokenInfo),
			zap.Any("request", sendCoinRequest))
		api_handler.InternalError(w, "internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}
