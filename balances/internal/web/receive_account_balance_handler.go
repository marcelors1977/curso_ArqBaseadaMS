package web

import (
	"encoding/json"
	"net/http"

	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/usercase/receive_balance"
)

type ReceiveWebAccountBalanceHandler struct {
	ReceiveAccountBalanceUseCase receive_balance.CreateAccountBalanceUseCase
}

func NewReceiveWebAccountBalanceHandler(receiveAcctBalUseCase receive_balance.CreateAccountBalanceUseCase) *ReceiveWebAccountBalanceHandler {
	return &ReceiveWebAccountBalanceHandler{
		ReceiveAccountBalanceUseCase: receiveAcctBalUseCase,
	}
}

func (h *ReceiveWebAccountBalanceHandler) ReceiveAccountBalance(w http.ResponseWriter, r *http.Request) {
	var dto receive_balance.CreateAccountBalanceInputDto

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err = h.ReceiveAccountBalanceUseCase.Execute(ctx, dto)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
