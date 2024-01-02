package gateway

import "github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"

type AccountGateway interface {
	Save(account *entity.Account) error
	FindById(id string) (*entity.Account, error)
	UpdateBalance(account *entity.Account) error
}
