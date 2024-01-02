package gateway

import "github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}
