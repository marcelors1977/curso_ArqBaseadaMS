package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/usercase/create_account"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/usercase/update_account"
)

type WebAccountHandler struct {
	CreateAccountUseCase create_account.CreateAccountUseCase
	UpdateAccountUseCase update_account.UpdateAccountUseCase
}

func NewWebAccountHandler(createAccountUseCase create_account.CreateAccountUseCase, updateAccountUseCase update_account.UpdateAccountUseCase) *WebAccountHandler {
	return &WebAccountHandler{
		CreateAccountUseCase: createAccountUseCase,
		UpdateAccountUseCase: updateAccountUseCase,
	}
}

func (h *WebAccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var dto create_account.CreateAccountInputDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := h.CreateAccountUseCase.Execute(dto)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
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

func (h *WebAccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var dto update_account.UpdateAccountInputDto
	account_id := chi.URLParam(r, "account_id")

	dto.AccountID = account_id

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := h.UpdateAccountUseCase.Execute(dto)
	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
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
