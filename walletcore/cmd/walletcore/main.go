package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/generator"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/database"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/event"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/event/handler"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/usercase/create_account"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/usercase/create_client"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/usercase/create_transaction"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/usercase/update_account"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/web"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/web/webserver"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/pkg/events"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/pkg/kafka"
	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/pkg/uow"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/wallet_core?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDb(db)
	accountDb := database.NewAccountDb(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDb", func(tx *sql.Tx) interface{} {
		return database.NewAccountDb(db)
	})

	uow.Register("TransactionDb", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDb(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	updateAccountUseCase := update_account.NewUpdateAccountUseCase(accountDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

	runMigration(db)
	runSeeders(db, *createTransactionUseCase)

	webserver := webserver.NewWebServer(":8080")

	clientHanlder := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase, *updateAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/client", clientHanlder.CreateClient)
	webserver.AddHandler("/account", accountHandler.CreateAccount)
	webserver.AddHandler("/account/update/{account_id}", accountHandler.UpdateAccount)
	webserver.AddHandler("/transaction", transactionHandler.CreateTransaction)

	fmt.Println("Rodando...")
	webserver.Start()

}

func runMigration(db *sql.DB) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		fmt.Println("error 1", err)
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file:///app/migrations", "mysql", driver)
	if err != nil {
		fmt.Println("error 2", err)
		panic(err)
	}

	m.Down()
	m.Up()
}

func runSeeders(db *sql.DB, createTransactionUseCase create_transaction.CreateTransactionUseCase) {
	a, err := generator.GenerateClients()
	if err != nil {
		fmt.Println(err)
	}

	clientDB := database.NewClientDb(db)
	for _, client := range a {
		err := clientDB.Save(&client)
		if err != nil {
			fmt.Println(err)
		}
	}

	b, err := generator.GenerateAccounts(a)
	if err != nil {
		fmt.Println(err)
	}

	accountDB := database.NewAccountDb(db)
	for _, account := range b {
		err := accountDB.Save(&account)
		if err != nil {
			fmt.Println(err)
		}
	}

	for i := 0; i < 30; i++ {
		account1 := b[rand.Intn(len(b))]
		account2 := b[rand.Intn(len(b))]

		transactionInputDto := create_transaction.CreateTransactionInputDto{
			AccountIDFrom: account1.ID,
			AccountIDTo:   account2.ID,
			Amount:        100,
		}

		ctx := context.Background()

		_, err = createTransactionUseCase.Execute(ctx, transactionInputDto)
		if err != nil {
			fmt.Println(err)
		}
	}
}
