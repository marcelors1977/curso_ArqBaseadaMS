package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	transactionDb *TransactionDb
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec(`CREATE TABLE transactions (account_id_from varchar(255), account_id_to varchar(255), amount float, date_transaction datetime)`)
	s.transactionDb = NewTransactionDb(db)
}

func (s *TransactionDBTestSuite) TearDownTest() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestSave() {
	transaction, _ := entity.NewTransaction("123", "321", 10, time.Now())
	err := s.transactionDb.Save(transaction)
	s.Nil(err)
}
