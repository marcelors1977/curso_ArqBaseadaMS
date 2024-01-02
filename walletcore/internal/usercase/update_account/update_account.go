package update_account

import (
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/gateway"
)

type UpdateAccountInputDto struct {
	AccountID string  `json:"id"`
	Balance   float64 `json:"balance"`
}

type UpdateAccountOutputDto struct {
	ID       string
	ClientID string
	Balance  float64
}

type UpdateAccountUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewUpdateAccountUseCase(a gateway.AccountGateway) *UpdateAccountUseCase {
	return &UpdateAccountUseCase{
		AccountGateway: a,
	}
}

func (uc *UpdateAccountUseCase) Execute(input UpdateAccountInputDto) (*UpdateAccountOutputDto, error) {
	account, err := uc.AccountGateway.FindById(input.AccountID)
	if err != nil {
		return nil, err
	}

	if input.Balance >= 0 {
		account.Credit((input.Balance))
	} else {
		account.Debit(input.Balance * -1)
	}

	err = uc.AccountGateway.UpdateBalance(account)
	if err != nil {
		return nil, err
	}

	return &UpdateAccountOutputDto{
		ID:       account.ID,
		ClientID: account.Client.ID,
		Balance:  account.Balance,
	}, nil
}
