package mocks

import (
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Save(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}
