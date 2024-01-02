package receive_trasaction

import (
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/gateway"
)

type ReceiveTransactionInputDto struct {
	AccountIdFrom   string
	AccountIdTo     string
	Amount          float64
	DateTransaction time.Time
}

type ReceiveTransactionUseCase struct {
	transactionGateway gateway.TransactionGateway
}

func NewReceiveTransactionUseCase(transactionGateway gateway.TransactionGateway) *ReceiveTransactionUseCase {
	return &ReceiveTransactionUseCase{
		transactionGateway: transactionGateway,
	}
}

func (uc *ReceiveTransactionUseCase) Execute(input ReceiveTransactionInputDto) error {
	transaction, err := entity.NewTransaction(input.AccountIdFrom, input.AccountIdTo, input.Amount, input.DateTransaction)
	if err != nil {
		return err
	}

	err = uc.transactionGateway.Save(transaction)
	if err != nil {
		return err
	}

	return nil
}
