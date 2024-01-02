package entity

import (
	"errors"
	"time"
)

type Transaction struct {
	AccountIdFrom   string    `json:"account_id_from"`
	AccountIdTo     string    `json:"account_id_to"`
	Amount          float64   `json:"amount"`
	DateTransaction time.Time `json:"date_transaction"`
}

func NewTransaction(_accountIdFrom string, _accountIdTo string, _amount float64, _dateTransaction time.Time) (*Transaction, error) {
	transaction := &Transaction{
		AccountIdFrom:   _accountIdFrom,
		AccountIdTo:     _accountIdTo,
		Amount:          _amount,
		DateTransaction: _dateTransaction,
	}

	err := transaction.Validate()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (a *Transaction) Validate() error {
	if a.AccountIdFrom == "" {
		return errors.New("AccountIdFrom id is required")
	}

	if a.AccountIdTo == "" {
		return errors.New("AccountIdTo id is required")
	}

	if a.DateTransaction.IsZero() {
		return errors.New("date transaction is required")
	}

	return nil
}
