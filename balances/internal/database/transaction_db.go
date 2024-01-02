package database

import (
	"database/sql"

	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
)

type TransactionDb struct {
	DB *sql.DB
}

func NewTransactionDb(db *sql.DB) *TransactionDb {
	return &TransactionDb{
		DB: db,
	}
}

func (a *TransactionDb) Save(transaction *entity.Transaction) error {
	stmt, err := a.DB.Prepare("INSERT INTO transactions (account_id_from, account_id_to, amount, date_transaction) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(transaction.AccountIdFrom, transaction.AccountIdTo, transaction.Amount, transaction.DateTransaction)
	if err != nil {
		return err
	}

	return nil
}
