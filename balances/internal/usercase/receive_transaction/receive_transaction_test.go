package receive_trasaction

import (
	"testing"
	"time"

	mocks "github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/usercase/_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReceiveTransaction_Save(t *testing.T) {
	mockGateway := &mocks.TransactionGatewayMock{}

	uc := NewReceiveTransactionUseCase(mockGateway)

	mockGateway.On("Save", mock.Anything).Return(nil)
	inputDTO := ReceiveTransactionInputDto{
		AccountIdFrom:   "123",
		AccountIdTo:     "456",
		Amount:          100,
		DateTransaction: time.Now(),
	}

	err := uc.Execute(inputDTO)

	assert.Nil(t, err)
	mockGateway.AssertExpectations(t)
	mockGateway.AssertNumberOfCalls(t, "Save", 1)

}
