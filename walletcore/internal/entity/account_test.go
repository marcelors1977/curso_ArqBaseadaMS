package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	client, _ := NewClient("John Doe", "john@doe")
	account := NewAccount(client)
	assert.NotNil(t, account)
	assert.Equal(t, client.ID, account.Client.ID)
	assert.Equal(t, client.Name, account.Client.Name)
	assert.Equal(t, client.Email, account.Client.Email)
	assert.Equal(t, float64(0), account.Balance)
	assert.Equal(t, account.UpdatedAt, time.Time{})
}

func TestCreateAccountWithNilClient(t *testing.T) {
	account := NewAccount(nil)
	assert.Nil(t, account)
}

func TestCreditAccount(t *testing.T) {
	client, _ := NewClient("John Doe", "john@doe")
	account := NewAccount(client)
	account.Credit(100)
	assert.Equal(t, float64(100), account.Balance)
	assert.NotEqual(t, time.Time{}, account.UpdatedAt)
}

func TestDebitAccount(t *testing.T) {
	client, _ := NewClient("John Doe", "john@doe")
	account := NewAccount(client)
	account.Credit(100)
	account.Debit(50)
	assert.Equal(t, float64(50), account.Balance)
	assert.NotEqual(t, time.Time{}, account.UpdatedAt)
}
