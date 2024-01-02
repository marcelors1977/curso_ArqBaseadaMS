package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/walletcore/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type ClientDBTestSuite struct {
	suite.Suite
	db       *sql.DB
	clientDb *ClientDb
}

func (s *ClientDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec(`CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at datetime, updated_at datetime)`)
	s.clientDb = NewClientDb(db)
}

func (s *ClientDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}

func (s *ClientDBTestSuite) TestGet() {
	client, _ := entity.NewClient("John Doe", "john@doe")
	s.clientDb.Save(client)

	clientDB, err := s.clientDb.Get(client.ID)
	s.Nil(err)
	s.Equal(client.ID, clientDB.ID)
	s.Equal(client.Name, clientDB.Name)
	s.Equal(client.Email, clientDB.Email)
	s.Equal(clientDB.UpdatedAt, time.Time{})
}

func (s *ClientDBTestSuite) TestSave() {
	client := &entity.Client{
		ID:    "123",
		Name:  "John Doe",
		Email: "john@doe",
	}

	err := s.clientDb.Save(client)
	s.Nil(err)
}
