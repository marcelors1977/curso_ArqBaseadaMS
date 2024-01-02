package database

import (
	"database/sql"
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"
)

type ClientDb struct {
	DB *sql.DB
}

func NewClientDb(db *sql.DB) *ClientDb {
	return &ClientDb{
		DB: db,
	}
}

func (c *ClientDb) Get(id string) (*entity.Client, error) {
	client := &entity.Client{}
	var createdAt sql.NullTime
	var updatedAt sql.NullTime

	stmt, err := c.DB.Prepare(`SELECT id, name, email, created_at, updated_at FROM clients WHERE id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	if err := row.Scan(&client.ID, &client.Name, &client.Email, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	if createdAt.Valid {
		client.CreatedAt = createdAt.Time
	} else {
		client.CreatedAt = time.Time{}
	}

	if updatedAt.Valid {
		client.UpdatedAt = updatedAt.Time
	} else {
		client.UpdatedAt = time.Time{}
	}

	return client, nil
}

func (c *ClientDb) Save(client *entity.Client) error {
	stmt, err := c.DB.Prepare("INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(client.ID, client.Name, client.Email, client.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
