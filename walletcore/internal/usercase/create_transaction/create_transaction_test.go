package create_transaction

import (
	"context"
	"testing"

	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/event"
	mocks "github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/usercase/_mocks"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	client1, _ := entity.NewClient("John Doe", "john@doe")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("Jane Doe", "jane@doe")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	inputDTO := CreateTransactionInputDto{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}

	dispatcher := events.NewEventDispatcher()
	eventTransaction := event.NewTransactionCreated()
	eventBalance := event.NewBalanceUpdated()
	ctx := context.Background()

	uc := NewCreateTransactionUseCase(mockUow, dispatcher, eventTransaction, eventBalance)
	output, err := uc.Execute(ctx, inputDTO)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
