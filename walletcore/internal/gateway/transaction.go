package gateway

import "github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
