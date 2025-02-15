package send_coin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/inna-maikut/avito-shop/internal/api"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-shop/internal/model"
)

type Handler struct {
	coinSending coinSending
}

func New(coinSending coinSending) (*Handler, error) {
	return &Handler{
		coinSending: coinSending,
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

		fmt.Println(err)
		api_handler.InternalError(w, "internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}
