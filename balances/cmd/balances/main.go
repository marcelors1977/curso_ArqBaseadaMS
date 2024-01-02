package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/database"
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/usercase/get_account_balance"
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/usercase/receive_balance"
	receive_trasaction "github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/usercase/receive_transaction"
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/web"
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/web/webserver"
	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/pkg/uow"
)

type BalancesTopic struct {
	Name    string
	Payload struct {
		AccountID       string
		AccountBalance  float64
		DateTransaction time.Time
	}
}

type TransactionsTopic struct {
	Name    string
	Payload struct {
		ID              string
		AccountIDFrom   string
		AccountIDTo     string
		Amount          float64
		DateTransaction time.Time
	}
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/balances?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer db.Close()
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	uow := uow.NewUow(db)

	uow.Register("AccountBalanceDb", func(tx *sql.Tx) interface{} {
		return database.NewAccountBalanceDb(db)
	})

	m, err := migrate.NewWithDatabaseInstance("file:///app/migrations", "mysql", driver)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	m.Down()
	m.Up()

	httpStart(uow, db)

	consumerStart(uow, db)
}

func httpStart(uow uow.UowInterface, db *sql.DB) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true

	go func() {
		webStartStatus := false

		for run {
			select {
			case sig := <-sigchan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				run = false
			default:
				if !webStartStatus {
					accountBalDB := database.NewAccountBalanceDb(db)
					transactionDB := database.NewTransactionDb(db)

					getAcctBalUseCase := get_account_balance.NewGetAccountBalanceUseCase(accountBalDB)
					receiveAcctBalUseCase := receive_balance.NewCreateAccountBalanceUseCase(uow)
					receiveTransactionUseCase := receive_trasaction.NewReceiveTransactionUseCase(transactionDB)

					webserver := webserver.NewWebServer(":3003")

					getBalanceHanlder := web.NewGetWebAccountBalanceHandler(*getAcctBalUseCase)
					receiveAcctBalHandler := web.NewReceiveWebAccountBalanceHandler(*receiveAcctBalUseCase)
					receiveTransactionHandler := web.NewReceiveWebTransactionHandler(*receiveTransactionUseCase)

					webserver.AddHandler("/balances/search/{account_id}", getBalanceHanlder.GetAccountBalance)
					webserver.AddHandler("/balances", receiveAcctBalHandler.ReceiveAccountBalance)
					webserver.AddHandler("/transaction", receiveTransactionHandler.ReceiveWebTransaction)

					fmt.Println("Webserver is running...")
					webStartStatus = true

					webserver.Start()
				}
			}
		}
	}()

}

func consumerStart(uow uow.UowInterface, db *sql.DB) {

	consumer, err := ckafka.NewConsumer(&ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{"balances", "transactions"}, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	consumeBalanceTopic := &BalancesTopic{}
	consumeTransactionTopic := &TransactionsTopic{}
	transactionDb := database.NewTransactionDb(db)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	fmt.Println("Consumer is running...")

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := consumer.ReadMessage(time.Duration(10) * time.Second)
			if err != nil {
				continue
			}

			switch string(*msg.TopicPartition.Topic) {
			case "balances":
				err := consumeBalanceTopic.ConsumeBalanceMsg(msg.Value, uow)
				if err != nil {
					fmt.Println(err)
					continue
				}

			case "transactions":
				err := consumeTransactionTopic.ConsumeTransactionMsg(msg.Value, transactionDb)
				if err != nil {
					fmt.Println(err)
					continue
				}

			default:
				fmt.Println("Consuming a unknown topic...", string(msg.Key))
			}
		}
	}
}

func (m *BalancesTopic) ConsumeBalanceMsg(msg []byte, uow uow.UowInterface) error {
	uc_createAccountBalance := receive_balance.NewCreateAccountBalanceUseCase(uow)

	ctx := context.Background()

	err := json.Unmarshal(msg, m)
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", m.Payload)

	inputDto := receive_balance.CreateAccountBalanceInputDto{
		AccountId:       m.Payload.AccountID,
		CurrentBalance:  m.Payload.AccountBalance,
		DateTransaction: m.Payload.DateTransaction,
	}

	err = uc_createAccountBalance.Execute(ctx, inputDto)
	if err != nil {
		return err
	}

	return nil
}

func (m *TransactionsTopic) ConsumeTransactionMsg(msg []byte, transactionDb *database.TransactionDb) error {
	uc_receiveTransaction := receive_trasaction.NewReceiveTransactionUseCase(transactionDb)

	err := json.Unmarshal(msg, m)
	if err != nil {
		return err
	}

	inputDTO := receive_trasaction.ReceiveTransactionInputDto{
		AccountIdFrom:   m.Payload.AccountIDFrom,
		AccountIdTo:     m.Payload.AccountIDTo,
		Amount:          m.Payload.Amount,
		DateTransaction: m.Payload.DateTransaction,
	}

	fmt.Printf("%#v\n", m.Payload)

	err = uc_receiveTransaction.Execute(inputDTO)
	if err != nil {
		return err
	}

	return nil
}
