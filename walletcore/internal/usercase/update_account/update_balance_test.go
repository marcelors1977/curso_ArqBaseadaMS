package update_account

import (
	"testing"

	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"
	mocks "github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/usercase/_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateAccountUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "john@doe")
	account1 := entity.NewAccount(client1)

	accountMock := &mocks.AccountGatewayMock{}
	accountMock.On("FindById", mock.Anything).Return(account1, nil)
	accountMock.On("UpdateBalance", mock.Anything).Return(nil)

	uc := NewUpdateAccountUseCase(accountMock)
	inputDTO := UpdateAccountInputDto{
		AccountID: account1.ID,
	}

	output, err := uc.Execute(inputDTO)
	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	accountMock.AssertExpectations(t)
	accountMock.AssertNumberOfCalls(t, "UpdateBalance", 1)
	assert.Equal(t, account1.Balance, output.Balance)
}
