package handler

import (
	"fmt"
	"sync"

	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/pkg/events"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/pkg/kafka"
)

type UpdateBalanceKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewUpdateBalanceKafkaHandler(kafka *kafka.Producer) *UpdateBalanceKafkaHandler {
	return &UpdateBalanceKafkaHandler{
		Kafka: kafka,
	}
}

func (h *UpdateBalanceKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()

	h.Kafka.Publish(message, nil, "balances")
	fmt.Println("UpdateBalanceKafkaHandler", message.GetPayload())
}
