package entity

import (
	"errors"
	"time"
)

type AccountBalance struct {
	AccountId       string
	CurrentBalance  float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DateTransaction time.Time
}

func NewAccountBalance(_accountId string) (*AccountBalance, error) {
	accountBalance := &AccountBalance{
		AccountId:       _accountId,
		CurrentBalance:  0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Time{},
		DateTransaction: time.Now(),
	}

	err := accountBalance.Validate()
	if err != nil {
		return nil, err
	}

	return accountBalance, nil
}

func (a *AccountBalance) Validate() error {
	if a.AccountId == "" {
		return errors.New("account id is required")
	}

	if a.DateTransaction.IsZero() {
		return errors.New("date transaction is required")
	}

	return nil
}

func (a *AccountBalance) UpdateAccountBalance(_currentBalance float64, _date_transaction time.Time) error {
	a.CurrentBalance = _currentBalance
	a.UpdatedAt = time.Now().Local()
	a.DateTransaction = _date_transaction

	err := a.Validate()
	if err != nil {
		return err
	}

	return nil
}
