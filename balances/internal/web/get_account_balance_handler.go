package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/usercase/get_account_balance"
)

type GetWebAccountBalanceHandler struct {
	GetAccountBalanceUseCase get_account_balance.GetAccountBalanceUseCase
}

func NewGetWebAccountBalanceHandler(getAcctBalUseCase get_account_balance.GetAccountBalanceUseCase) *GetWebAccountBalanceHandler {
	return &GetWebAccountBalanceHandler{
		GetAccountBalanceUseCase: getAcctBalUseCase,
	}
}

func (h *GetWebAccountBalanceHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	var dto get_account_balance.GetAccountBalanceInputDto
	account_id := chi.URLParam(r, "account_id")

	dto.AccountId = account_id

	output, err := h.GetAccountBalanceUseCase.GetAccountBalance(dto)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
