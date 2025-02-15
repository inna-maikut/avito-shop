package info

import (
	"net/http"

	"github.com/inna-maikut/avito-shop/internal/api"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/api_handler"
	"github.com/inna-maikut/avito-shop/internal/infrastructure/jwt"
	"github.com/inna-maikut/avito-shop/internal/model"
)

type Handler struct {
	infoCollecting infoCollecting
}

func New(infoCollecting infoCollecting) (*Handler, error) {
	return &Handler{
		infoCollecting: infoCollecting,
	}, nil
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenInfo := jwt.TokenInfoFromContext(r.Context())

	info, err := h.infoCollecting.Collect(ctx, tokenInfo.EmployeeID)
	if err != nil {
		api_handler.InternalError(w, "internal server error")
		return
	}

	api_handler.OK(w, convertToResponse(info))
}

func convertToResponse(info model.EmployeeInfo) api.InfoResponse {
	type apiInventoryItem = struct {
		Quantity *int    `json:"quantity,omitempty"`
		Type     *string `json:"type,omitempty"`
	}
	inventory := make([]apiInventoryItem, 0, len(info.Inventory))
	for _, i := range info.Inventory {
		inventory = append(inventory, apiInventoryItem{
			Quantity: pointerOfInt(i.Quantity),
			Type:     &i.MerchName,
		})
	}

	type apiReceivedTransaction = struct {
		Amount   *int    `json:"amount,omitempty"`
		FromUser *string `json:"fromUser,omitempty"`
	}
	received := make([]apiReceivedTransaction, 0, len(info.ReceivedTransactions))
	for _, t := range info.ReceivedTransactions {
		received = append(received, apiReceivedTransaction{
			Amount:   pointerOfInt(t.Amount),
			FromUser: pointerOf(t.CounterpartyUsername),
		})
	}

	type apiSentTransaction = struct {
		Amount *int    `json:"amount,omitempty"`
		ToUser *string `json:"toUser,omitempty"`
	}
	sent := make([]apiSentTransaction, 0, len(info.SentTransactions))
	for _, t := range info.SentTransactions {
		sent = append(sent, apiSentTransaction{
			Amount: pointerOfInt(t.Amount),
			ToUser: pointerOf(t.CounterpartyUsername),
		})
	}

	return api.InfoResponse{
		Coins: pointerOfInt(info.Coins),
		CoinHistory: &struct {
			Received *[]apiReceivedTransaction `json:"received,omitempty"`
			Sent     *[]apiSentTransaction     `json:"sent,omitempty"`
		}{
			Received: &received,
			Sent:     &sent,
		},
		Inventory: &inventory,
	}
}

func pointerOf[T any](v T) *T {
	return &v
}

func pointerOfInt[T int64 | int32 | int](v T) *int {
	return pointerOf(int(v))
}
