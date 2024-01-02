package database

import (
	"database/sql"
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
)

type AccountBalanceDb struct {
	DB *sql.DB
}

func NewAccountBalanceDb(db *sql.DB) *AccountBalanceDb {
	return &AccountBalanceDb{
		DB: db,
	}
}

func (a *AccountBalanceDb) Get(accountId string) (*entity.AccountBalance, error) {
	var UpdatedAt sql.NullTime

	accountBalance := &entity.AccountBalance{}
	stmt, err := a.DB.Prepare("SELECT account_id, current_balance, created_at, updated_at FROM account_balances WHERE account_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(accountId)
	if err := row.Scan(&accountBalance.AccountId, &accountBalance.CurrentBalance, &accountBalance.CreatedAt, &UpdatedAt); err != nil {
		return nil, err
	}

	if UpdatedAt.Valid {
		accountBalance.UpdatedAt = UpdatedAt.Time
	} else {
		accountBalance.UpdatedAt = time.Time{}
	}

	return accountBalance, nil
}

func (a *AccountBalanceDb) Save(accountBalance *entity.AccountBalance) error {
	stmt, err := a.DB.Prepare("INSERT INTO account_balances (account_id, current_balance, created_at, updated_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(accountBalance.AccountId, accountBalance.CurrentBalance, accountBalance.CreatedAt, accountBalance.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountBalanceDb) Update(accountBalance *entity.AccountBalance) error {
	stmt, err := a.DB.Prepare("UPDATE account_balances SET current_balance = ?, updated_at = ? WHERE account_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(accountBalance.CurrentBalance, accountBalance.UpdatedAt, accountBalance.AccountId)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountBalanceDb) SaveHistory(accountBalance *entity.AccountBalance) error {
	stmt, err := a.DB.Prepare("INSERT INTO account_balances_history (account_id, current_balance, date_transaction) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(accountBalance.AccountId, accountBalance.CurrentBalance, accountBalance.DateTransaction)
	if err != nil {
		return err
	}

	return nil
}
