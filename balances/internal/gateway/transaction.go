package gateway

import "github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"

type TransactionGateway interface {
	Save(transaction *entity.Transaction) error
}
