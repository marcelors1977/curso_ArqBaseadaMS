package web

import (
	"encoding/json"
	"net/http"

	receive_trasaction "github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/usercase/receive_transaction"
)

type ReceiveWebTransactionHandler struct {
	ReceiveTransactionUseCase receive_trasaction.ReceiveTransactionUseCase
}

func NewReceiveWebTransactionHandler(receiveTransactionUseCase receive_trasaction.ReceiveTransactionUseCase) *ReceiveWebTransactionHandler {
	return &ReceiveWebTransactionHandler{
		ReceiveTransactionUseCase: receiveTransactionUseCase,
	}
}

func (h *ReceiveWebTransactionHandler) ReceiveWebTransaction(w http.ResponseWriter, r *http.Request) {
	var dto receive_trasaction.ReceiveTransactionInputDto

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ReceiveTransactionUseCase.Execute(dto)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
