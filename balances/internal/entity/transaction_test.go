package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Validate(t *testing.T) {
	transaction := &Transaction{
		AccountIdFrom:   "",
		AccountIdTo:     "",
		Amount:          0,
		DateTransaction: time.Time{},
	}

	err := transaction.Validate()
	assert.Error(t, err, "AccountIdFrom id is required")

	transaction.AccountIdFrom = "123"

	err = transaction.Validate()
	assert.Error(t, err, "AccountIdTo id is required")

	transaction.AccountIdTo = "456"

	err = transaction.Validate()
	assert.Error(t, err, "date transaction is required")

	transaction.DateTransaction = time.Now()

	err = transaction.Validate()
	assert.Nil(t, err)
}

func TestTransaction_NewTransaction(t *testing.T) {
	transaction, err := NewTransaction("123", "456", 100, time.Now().Truncate(time.Minute))
	assert.Nil(t, err)
	assert.Equal(t, transaction.AccountIdFrom, "123")
	assert.Equal(t, transaction.AccountIdTo, "456")
	assert.Equal(t, transaction.Amount, float64(100))
	assert.Equal(t, transaction.DateTransaction, time.Now().Truncate(time.Minute))

	_, err = NewTransaction("", "456", 100, time.Now().Truncate(time.Minute))
	assert.NotNil(t, err)
	assert.Error(t, err, "AccountIdFrom id is required")
}
