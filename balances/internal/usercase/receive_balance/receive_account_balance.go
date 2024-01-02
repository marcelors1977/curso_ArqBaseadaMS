package receive_balance

import (
	"context"
	"database/sql"

	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/gateway"
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/pkg/uow"
)

type FnInterface interface {
	FnSave(ctx context.Context, uc *CreateAccountBalanceUseCase, acctBal *entity.AccountBalance) error
	FnUpdate(ctx context.Context, uc *CreateAccountBalanceUseCase, acctBal *entity.AccountBalance) error
}

type Fn struct{}

func (fn Fn) FnSave(ctx context.Context, uc *CreateAccountBalanceUseCase, acctBal *entity.AccountBalance) error {
	repo := uc.getAccountBalanceRepository(ctx)

	err := uc.Uow.Do(ctx, func(uow *uow.Uow) error {
		err := repo.Save(acctBal)
		if err != nil {
			return err
		}

		err = repo.SaveHistory(acctBal)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (fn Fn) FnUpdate(ctx context.Context, uc *CreateAccountBalanceUseCase, acctBal *entity.AccountBalance) error {
	repo := uc.getAccountBalanceRepository(ctx)

	err := uc.Uow.Do(ctx, func(uow *uow.Uow) error {
		err := repo.Update(acctBal)
		if err != nil {
			return err
		}

		err = repo.SaveHistory(acctBal)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

type CreateAccountBalanceInputDto struct {
	AccountId       string    `json:"AccountId"`
	CurrentBalance  float64   `json:"AccountBalance"`
	DateTransaction time.Time `json:"DateTransaction"`
}

type CreateAccountBalanceOutputDto struct{}

type CreateAccountBalanceUseCase struct {
	Uow uow.UowInterface
	Fn  FnInterface
}

func NewCreateAccountBalanceUseCase(Uow uow.UowInterface) *CreateAccountBalanceUseCase {
	return &CreateAccountBalanceUseCase{
		Uow: Uow,
		Fn:  Fn{},
	}
}

func (uc *CreateAccountBalanceUseCase) Execute(ctx context.Context, input CreateAccountBalanceInputDto) error {
	accountBalance, err := entity.NewAccountBalance(input.AccountId)
	if err != nil {
		return err
	}

	err = accountBalance.UpdateAccountBalance(input.CurrentBalance, input.DateTransaction)
	if err != nil {
		return err
	}

	repo := uc.getAccountBalanceRepository(ctx)
	accountBalanceDB, err := repo.Get(accountBalance.AccountId)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		err = uc.Fn.FnSave(ctx, uc, accountBalance)
		if err != nil {
			return err
		}
	} else {
		if accountBalanceDB.UpdatedAt.Before(accountBalance.UpdatedAt) {
			err = uc.Fn.FnUpdate(ctx, uc, accountBalance)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (uc *CreateAccountBalanceUseCase) getAccountBalanceRepository(ctx context.Context) gateway.AccountBalanceGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountBalanceDb")
	if err != nil {
		panic(err)
	}

	return repo.(gateway.AccountBalanceGateway)
}
