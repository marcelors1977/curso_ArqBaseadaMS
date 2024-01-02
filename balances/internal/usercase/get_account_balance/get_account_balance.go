package get_account_balance

import (
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/gateway"
)

type GetAccountBalanceInputDto struct {
	AccountId string
}

type GetAccountBalanceOutputDto struct {
	AccountId string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetAccountBalanceUseCase struct {
	accountBalanceGateway gateway.AccountBalanceGateway
}

func NewGetAccountBalanceUseCase(accountBalanceGateway gateway.AccountBalanceGateway) *GetAccountBalanceUseCase {
	return &GetAccountBalanceUseCase{
		accountBalanceGateway: accountBalanceGateway,
	}
}

func (uc *GetAccountBalanceUseCase) GetAccountBalance(input GetAccountBalanceInputDto) (*GetAccountBalanceOutputDto, error) {
	accountBalance, err := uc.accountBalanceGateway.Get(input.AccountId)
	if err != nil {
		return nil, err
	}

	return &GetAccountBalanceOutputDto{
		AccountId: accountBalance.AccountId,
		Balance:   accountBalance.CurrentBalance,
		CreatedAt: accountBalance.CreatedAt,
		UpdatedAt: accountBalance.UpdatedAt,
	}, nil
}
