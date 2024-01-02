package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccountBalance(t *testing.T) {
	accountBalance, err := NewAccountBalance("123")
	assert.NotNil(t, accountBalance)
	assert.Nil(t, err)
	assert.Equal(t, "123", accountBalance.AccountId)
	assert.Equal(t, float64(0), accountBalance.CurrentBalance)
	assert.NotEqual(t, accountBalance.CreatedAt, time.Time{})
	assert.NotEqual(t, accountBalance.DateTransaction, time.Time{})
	assert.Equal(t, accountBalance.UpdatedAt, time.Time{})
}

func TestCreatingAccountBalanceWithInvalidArgs(t *testing.T) {
	accountBalance, err := NewAccountBalance("")
	assert.Nil(t, accountBalance)
	assert.NotNil(t, err)
	assert.Error(t, err, "account id is required")

	accountBalance1 := &AccountBalance{
		AccountId: "123",
	}

	err = accountBalance1.Validate()
	assert.NotNil(t, err)
	assert.Error(t, err, "date transaction is required")
}

func TestUpdateAccountBalance(t *testing.T) {
	accountBalance, _ := NewAccountBalance("123")
	assert.NotNil(t, accountBalance)
	assert.Equal(t, accountBalance.CurrentBalance, float64(0))

	dateTransaction := time.Now().Add(time.Minute * 5)
	accountBalance.UpdateAccountBalance(100, dateTransaction)
	assert.Equal(t, accountBalance.CurrentBalance, float64(100))
	assert.Contains(t, accountBalance.UpdatedAt.String(), time.Now().Format("2006-01-02 15:04:05"))
	assert.Contains(t, accountBalance.DateTransaction.String(), dateTransaction.Format("2006-01-02 15:04:05"))

	var accountBalance1 *AccountBalance
	assert.Nil(t, accountBalance1)
}
