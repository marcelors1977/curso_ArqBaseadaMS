package get_account_balance

import (
	"testing"

	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
	mocks "github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/usercase/_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAccountBalance(t *testing.T) {
	getAcctBalGateway := &mocks.AccountBalanceGatewayMock{}

	uc := NewGetAccountBalanceUseCase(getAcctBalGateway)

	accountBalance, err := entity.NewAccountBalance("123")

	getAcctBalGateway.On("Get", mock.Anything).Return(accountBalance, nil)

	assert.Nil(t, err)

	inputDTO := GetAccountBalanceInputDto{
		AccountId: accountBalance.AccountId,
	}

	output, err := uc.GetAccountBalance(inputDTO)

	assert.Nil(t, err)
	assert.Equal(t, accountBalance.AccountId, output.AccountId)
	getAcctBalGateway.AssertExpectations(t)
	getAcctBalGateway.AssertNumberOfCalls(t, "Get", 1)

}
