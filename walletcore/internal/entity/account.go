package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string  `faker:"uuid_digit"`
	Client    *Client `faker:"client"`
	Balance   float64 `faker:"oneof:600.0,700.0,800.0,900.0,1000.0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount(client *Client) *Account {
	if client == nil {
		return nil
	}

	account := &Account{
		ID:        uuid.New().String(),
		Client:    client,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}

	return account
}

func (a *Account) Credit(amount float64) {
	a.Balance += amount
	a.UpdatedAt = time.Now()
}

func (a *Account) Debit(amount float64) {
	a.Balance -= amount
	a.UpdatedAt = time.Now()
}
