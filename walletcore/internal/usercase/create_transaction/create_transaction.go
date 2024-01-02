package create_transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/gateway"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/pkg/events"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/pkg/uow"
)

type CreateTransactionInputDto struct {
	AccountIDFrom string  `json:"account_id_from"`
	AccountIDTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDto struct {
	ID              string
	AccountIDFrom   string
	AccountIDTo     string
	Amount          float64
	DateTransaction time.Time
}

type BalanceUpdatedOutputDto struct {
	AccountID       string
	AccountBalance  float64
	DateTransaction time.Time
}

type CreateTransactionUseCase struct {
	Uow                uow.UowInterface
	eventDispatcher    events.EventDispatcherInterface
	transactionCreated events.EventInterface
	balanceUpdated     events.EventInterface
}

func NewCreateTransactionUseCase(
	Uow uow.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	balanceUpdated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		Uow:                Uow,
		eventDispatcher:    eventDispatcher,
		transactionCreated: transactionCreated,
		balanceUpdated:     balanceUpdated,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDto) (*CreateTransactionOutputDto, error) {
	output := &CreateTransactionOutputDto{}
	balanceUpdated := make([]BalanceUpdatedOutputDto, 2)
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := uc.getAccountRepository(ctx)
		transactionRepository := uc.getTransactionRepository(ctx)

		accountFrom, err := accountRepository.FindById(input.AccountIDFrom)
		if err != nil {
			return err
		}
		accountTo, err := accountRepository.FindById(input.AccountIDTo)
		if err != nil {
			return err
		}
		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			return err
		}

		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}
		err = accountRepository.UpdateBalance(accountTo)
		if err != nil {
			return err
		}

		err = transactionRepository.Create(transaction)
		if err != nil {
			return err
		}

		output.ID = transaction.ID
		output.AccountIDFrom = input.AccountIDFrom
		output.AccountIDTo = input.AccountIDTo
		output.Amount = input.Amount
		if (transaction.UpdateAt != time.Time{}) {
			output.DateTransaction = transaction.UpdateAt.Truncate(time.Second)
			balanceUpdated[0].DateTransaction = transaction.UpdateAt.Truncate(time.Second)
			balanceUpdated[1].DateTransaction = transaction.UpdateAt.Truncate(time.Second)
		} else {
			output.DateTransaction = transaction.CreatedAt.Truncate(time.Second)
			balanceUpdated[0].DateTransaction = transaction.CreatedAt.Truncate(time.Second)
			balanceUpdated[1].DateTransaction = transaction.CreatedAt.Truncate(time.Second)
		}

		balanceUpdated[0].AccountID = input.AccountIDFrom
		balanceUpdated[0].AccountBalance = accountFrom.Balance
		balanceUpdated[1].AccountID = input.AccountIDTo
		balanceUpdated[1].AccountBalance = accountTo.Balance

		return nil
	})
	if err != nil {
		return nil, err
	}

	uc.transactionCreated.SetPayload(output)
	uc.eventDispatcher.Dispatch(uc.transactionCreated)

	for _, balance := range balanceUpdated {
		uc.balanceUpdated.SetPayload(balance)
		fmt.Printf("balanceUpdated: %#v\n", balance)
		uc.eventDispatcher.Dispatch(uc.balanceUpdated)
	}

	return output, nil
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDb")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := uc.Uow.GetRepository(ctx, "TransactionDb")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGateway)
}
