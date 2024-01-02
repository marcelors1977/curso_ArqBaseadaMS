package generator

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"
)

func GenerateClients() ([]entity.Client, error) {
	CustomerGeneratorDate()
	clients := []entity.Client{}
	err := faker.FakeData(
		&clients,
		options.WithRandomMapAndSliceMaxSize(100),
		options.WithRandomMapAndSliceMinSize(50),
	)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func CustomerGeneratorDate() {
	_ = faker.AddProvider("dateTimeNow", func(v reflect.Value) (interface{}, error) {
		return time.Now(), nil
	})
}

func GenerateAccounts(listClients []entity.Client) ([]entity.Account, error) {
	CustomGeneratorClientList(listClients)
	CustomerGeneratorDate()
	accounts := []entity.Account{}
	err := faker.FakeData(
		&accounts,
		options.WithRandomMapAndSliceMaxSize(150),
		options.WithRandomMapAndSliceMinSize(120),
	)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func CustomGeneratorClientList(listClients []entity.Client) {
	_ = faker.AddProvider("client", func(v reflect.Value) (interface{}, error) {
		obj := &listClients[rand.Intn(len(listClients))]
		return obj, nil
	})
}
