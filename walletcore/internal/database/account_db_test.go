package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDb *AccountDb
	client    *entity.Client
}

func (s *AccountDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec(`CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at datetime, updated_at datetime)`)
	db.Exec(`CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at datetime, updated_at datetime)`)
	s.accountDb = NewAccountDb(db)
	s.client, _ = entity.NewClient("John Doe", "john@doe")
	s.db.Exec("Insert into clients (id, name, email, created_at) values (?, ?, ?, ?)",
		s.client.ID, s.client.Name, s.client.Email, s.client.CreatedAt,
	)
}

func (s *AccountDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (s *AccountDBTestSuite) TestSave() {
	account := entity.NewAccount(s.client)
	account.Credit(1000.0)
	err := s.accountDb.Save(account)
	s.Nil(err)
}

func (s *AccountDBTestSuite) TestFindById() {
	account := entity.NewAccount(s.client)
	err := s.accountDb.Save(account)
	s.Nil(err)

	accountDB, err := s.accountDb.FindById(account.ID)
	s.Nil(err)
	s.Equal(account.ID, accountDB.ID)
	s.Equal(account.Balance, accountDB.Balance)
	s.Equal(account.Client.ID, accountDB.Client.ID)
	s.Equal(account.Client.Name, accountDB.Client.Name)
	s.Equal(account.Client.Email, accountDB.Client.Email)
	s.Equal(time.Time{}, accountDB.UpdatedAt)
}

func (s *AccountDBTestSuite) TestUpdateBalance() {
	account := entity.NewAccount(s.client)
	err := s.accountDb.Save(account)
	s.Nil(err)
	s.Equal(0.0, account.Balance)

	account.Credit(1000.0)
	err = s.accountDb.UpdateBalance(account)
	s.Nil(err)
	s.Equal(1000.0, account.Balance)
	accountFromDB, err := s.accountDb.FindById(account.ID)
	s.Nil(err)

	s.Equal(1000.0, accountFromDB.Balance)
	s.NotEqual(time.Time{}, accountFromDB.UpdatedAt)
	s.Equal(accountFromDB.UpdatedAt.Local(), account.UpdatedAt.Local())
}
