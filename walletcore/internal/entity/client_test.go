package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	client, err := NewClient("John Doe", "john@doe")
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "John Doe", client.Name)
	assert.Equal(t, "john@doe", client.Email)
	assert.Equal(t, client.UpdatedAt, time.Time{})
}

func TestCreateNewClientInvalidArgs(t *testing.T) {
	client, err := NewClient("", "")
	assert.NotNil(t, err)
	assert.Nil(t, client)
}

func TestUpdateClient(t *testing.T) {
	client, _ := NewClient("John Doe", "john@doe")
	err := client.Update("Jane Doe", "jane@doe")
	assert.Nil(t, err)
	assert.Equal(t, "Jane Doe", client.Name)
	assert.Equal(t, "jane@doe", client.Email)
	assert.NotEqual(t, client.UpdatedAt, time.Time{})
}

func TestUpdateClientWithInvalidArgs(t *testing.T) {
	client, _ := NewClient("John Doe", "john@doe")
	err := client.Update("", "john@doe")
	assert.Error(t, err, "name is required")
	assert.NotNil(t, err)
}

func TestAddAccountToClient(t *testing.T) {
	client, _ := NewClient("John Doe", "john@doe")
	account := NewAccount(client)
	err := client.AddAccount(account)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(client.Accounts))
}
