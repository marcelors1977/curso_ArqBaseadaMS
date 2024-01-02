package database

import (
	"database/sql"
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"
)

type AccountDb struct {
	DB *sql.DB
}

func NewAccountDb(db *sql.DB) *AccountDb {
	return &AccountDb{
		DB: db,
	}
}

func (a *AccountDb) FindById(id string) (*entity.Account, error) {
	var account entity.Account
	var client entity.Client
	account.Client = &client

	var accountUpdatedAt sql.NullTime
	var clientUpdatedAt sql.NullTime

	stmt, err := a.DB.Prepare("Select a.id, a.client_id, a.balance, a.created_at, a.updated_at, c.id, c.name, c.email, c.created_at, c.updated_at FROM accounts a INNER JOIN clients c ON a.client_id = c.id WHERE a.id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(
		&account.ID,
		&account.Client.ID,
		&account.Balance,
		&account.CreatedAt,
		&accountUpdatedAt,
		&client.ID,
		&client.Name,
		&client.Email,
		&client.CreatedAt,
		&clientUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if accountUpdatedAt.Valid {
		account.UpdatedAt = accountUpdatedAt.Time
	} else {
		account.UpdatedAt = time.Time{}
	}

	if clientUpdatedAt.Valid {
		client.UpdatedAt = clientUpdatedAt.Time
	} else {
		client.UpdatedAt = time.Time{}
	}

	return &account, nil
}

func (a *AccountDb) Save(account *entity.Account) error {
	stmt, err := a.DB.Prepare("INSERT INTO accounts (id, client_id, balance, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(account.ID, account.Client.ID, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountDb) UpdateBalance(account *entity.Account) error {
	stmt, err := a.DB.Prepare("UPDATE accounts SET balance = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.Balance, account.UpdatedAt, account.ID)
	if err != nil {
		return err
	}
	return nil
}
