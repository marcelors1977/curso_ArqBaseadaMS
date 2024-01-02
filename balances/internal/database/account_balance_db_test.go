package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type AccountBalanceDBTestSuite struct {
	suite.Suite
	db               *sql.DB
	accountBalanceDb *AccountBalanceDb
	accountBalance   *entity.AccountBalance
}

func (s *AccountBalanceDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec(`CREATE TABLE account_balances (account_id varchar(255), current_balance float, created_at datetime, updated_at datetime)`)
	db.Exec(`CREATE TABLE account_balances_history (account_id varchar(255), current_balance float, date_transaction datetime)`)
	s.accountBalanceDb = NewAccountBalanceDb(db)
	s.accountBalance, _ = entity.NewAccountBalance("123")
}

func (s *AccountBalanceDBTestSuite) TearDownTest() {
	defer s.SetupSuite()
	defer s.db.Close()
	s.db.Exec("DROP TABLE account_balances")
	s.db.Exec("DROP TABLE account_balances_history")
}

func TestAccountBalanceDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountBalanceDBTestSuite))
}

func (s *AccountBalanceDBTestSuite) TestGet() {
	s.accountBalanceDb.Save(s.accountBalance)

	accountBalanceDB, err := s.accountBalanceDb.Get(s.accountBalance.AccountId)
	s.Nil(err)
	s.Equal(s.accountBalance.AccountId, accountBalanceDB.AccountId)
	s.Equal(s.accountBalance.CurrentBalance, accountBalanceDB.CurrentBalance)
}

func (s *AccountBalanceDBTestSuite) TestSave() {
	err := s.accountBalanceDb.Save(s.accountBalance)
	s.Nil(err)
}

func (s *AccountBalanceDBTestSuite) TestUpdate() {
	_, err := s.accountBalanceDb.Get(s.accountBalance.AccountId)
	s.Equal(err, sql.ErrNoRows)

	s.accountBalanceDb.Save(s.accountBalance)

	accountBalanceDB, err := s.accountBalanceDb.Get(s.accountBalance.AccountId)
	s.Nil(err)
	s.Equal(accountBalanceDB.UpdatedAt, time.Time{})

	timeNow := time.Now()
	s.accountBalance.UpdatedAt = timeNow

	err = s.accountBalanceDb.Update(s.accountBalance)
	s.Nil(err)
	accountBalanceDB, err = s.accountBalanceDb.Get(s.accountBalance.AccountId)
	s.Nil(err)
	s.Equal(accountBalanceDB.UpdatedAt.Local(), timeNow.Local())
}

func (s *AccountBalanceDBTestSuite) TestSaveHistory() {
	err := s.accountBalanceDb.SaveHistory(s.accountBalance)
	s.Nil(err)
}
