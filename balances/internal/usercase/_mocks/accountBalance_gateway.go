package mocks

import (
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
	"github.com/stretchr/testify/mock"
)

type AccountBalanceGatewayMock struct {
	mock.Mock
}

func (m *AccountBalanceGatewayMock) Get(id string) (*entity.AccountBalance, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.AccountBalance), args.Error(1)
}

func (m *AccountBalanceGatewayMock) Save(account *entity.AccountBalance) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountBalanceGatewayMock) Update(account *entity.AccountBalance) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountBalanceGatewayMock) SaveHistory(account *entity.AccountBalance) error {
	args := m.Called(account)
	return args.Error(0)
}
