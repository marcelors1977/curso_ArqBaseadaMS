package gateway

import (
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
)

type AccountBalanceGateway interface {
	Get(id string) (*entity.AccountBalance, error)
	Save(account *entity.AccountBalance) error
	Update(account *entity.AccountBalance) error
	SaveHistory(account *entity.AccountBalance) error
}
