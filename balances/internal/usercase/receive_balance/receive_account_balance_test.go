package receive_balance

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/entity"
	mocks "github.com/marcelors1977/curso_ArqBaseadaMS/balances/internal/usercase/_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type FnMock struct {
	mock.Mock
}

func (m *FnMock) FnSave(ctx context.Context, uc *CreateAccountBalanceUseCase, acctBal *entity.AccountBalance) error {
	args := m.Called(ctx, uc, acctBal)
	return args.Error(0)
}

func (m *FnMock) FnUpdate(ctx context.Context, uc *CreateAccountBalanceUseCase, acctBal *entity.AccountBalance) error {
	args := m.Called(ctx, uc, acctBal)
	return args.Error(0)
}

type AcctBalTestSuite struct {
	suite.Suite
	ctx              context.Context
	mockUow          *mocks.UowMock
	gatewayMock      *mocks.AccountBalanceGatewayMock
	fnMock           *FnMock
	accountBalanceDB *entity.AccountBalance
	inputDTO         CreateAccountBalanceInputDto
	uc               CreateAccountBalanceUseCase
}

func (s *AcctBalTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.mockUow = &mocks.UowMock{}
	s.gatewayMock = &mocks.AccountBalanceGatewayMock{}
	s.fnMock = &FnMock{}

	s.accountBalanceDB, _ = entity.NewAccountBalance("1")

	s.inputDTO = CreateAccountBalanceInputDto{
		AccountId:       s.accountBalanceDB.AccountId,
		CurrentBalance:  s.accountBalanceDB.CurrentBalance,
		DateTransaction: s.accountBalanceDB.DateTransaction,
	}

	s.uc = CreateAccountBalanceUseCase{
		Uow: s.mockUow,
		Fn:  s.fnMock,
	}
}

func (s *AcctBalTestSuite) TearDownTest() {
	defer s.SetupSuite()
}

func TestAcctBalTestSuite(t *testing.T) {
	suite.Run(t, new(AcctBalTestSuite))
}

func (s *AcctBalTestSuite) TestReceiveAcctBalUseCase_SaveNewAcctBal() {
	s.fnMock.On("FnSave", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s.mockUow.On("GetRepository", mock.Anything, mock.Anything).Return(s.gatewayMock, nil)
	s.gatewayMock.On("Get", s.accountBalanceDB.AccountId).Return(s.accountBalanceDB, sql.ErrNoRows)

	err := s.uc.Execute(s.ctx, s.inputDTO)

	s.Nil(err)
	s.mockUow.AssertExpectations(s.T())
	s.mockUow.AssertNumberOfCalls(s.T(), "GetRepository", 1)
	s.gatewayMock.AssertExpectations(s.T())
	s.gatewayMock.AssertNumberOfCalls(s.T(), "Get", 1)
	s.fnMock.AssertExpectations(s.T())
	s.fnMock.AssertNumberOfCalls(s.T(), "FnSave", 1)
	s.fnMock.AssertNumberOfCalls(s.T(), "FnUpdate", 0)
}

func (s *AcctBalTestSuite) TestReceiveAcctBalUseCase_UpdtWhenExistsAcctBal() {

	s.accountBalanceDB.UpdatedAt = time.Time{}

	s.fnMock.On("FnUpdate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s.mockUow.On("GetRepository", mock.Anything, mock.Anything).Return(s.gatewayMock, nil)
	s.gatewayMock.On("Get", s.accountBalanceDB.AccountId).Return(s.accountBalanceDB, nil)

	err := s.uc.Execute(s.ctx, s.inputDTO)

	s.Nil(err)
	s.mockUow.AssertExpectations(s.T())
	s.mockUow.AssertNumberOfCalls(s.T(), "GetRepository", 1)
	s.gatewayMock.AssertExpectations(s.T())
	s.gatewayMock.AssertNumberOfCalls(s.T(), "Get", 1)
	s.fnMock.AssertExpectations(s.T())
	s.fnMock.AssertNumberOfCalls(s.T(), "FnSave", 0)
	s.fnMock.AssertNumberOfCalls(s.T(), "FnUpdate", 1)
}

func (s *AcctBalTestSuite) TestReceiveAcctBalUseCase_WhenErrorOnGetRepo() {
	s.mockUow.On("GetRepository", mock.Anything, mock.Anything).Return(s.gatewayMock, nil)
	s.gatewayMock.On("Get", s.accountBalanceDB.AccountId).Return(s.accountBalanceDB, errors.New("Unexpected error"))

	err := s.uc.Execute(s.ctx, s.inputDTO)

	s.NotNil(err)
	assert.NotEqual(s.T(), err, sql.ErrNoRows)
	s.fnMock.AssertNumberOfCalls(s.T(), "FnSave", 0)
	s.fnMock.AssertNumberOfCalls(s.T(), "FnUpdate", 0)
	s.mockUow.AssertNumberOfCalls(s.T(), "GetRepository", 1)
	s.gatewayMock.AssertNumberOfCalls(s.T(), "Get", 1)
}

func (s *AcctBalTestSuite) TestReceiveAcctBalUseCase_WhenNothingToUpdate() {
	s.accountBalanceDB.UpdatedAt = time.Now().Add(2 * time.Hour)

	s.mockUow.On("GetRepository", mock.Anything, mock.Anything).Return(s.gatewayMock, nil)
	s.gatewayMock.On("Get", s.accountBalanceDB.AccountId).Return(s.accountBalanceDB, nil)

	err := s.uc.Execute(s.ctx, s.inputDTO)

	s.Nil(err)
	s.mockUow.AssertExpectations(s.T())
	s.mockUow.AssertNumberOfCalls(s.T(), "GetRepository", 1)
	s.gatewayMock.AssertExpectations(s.T())
	s.gatewayMock.AssertNumberOfCalls(s.T(), "Get", 1)
	s.fnMock.AssertExpectations(s.T())
	s.fnMock.AssertNumberOfCalls(s.T(), "FnSave", 0)
	s.fnMock.AssertNumberOfCalls(s.T(), "FnUpdate", 0)
}

func (s *AcctBalTestSuite) TestReceiveAcctBalUseCase_ErrorFnSave() {
	s.mockUow.On("GetRepository", mock.Anything, mock.Anything).Return(s.gatewayMock, nil)
	s.gatewayMock.On("Get", s.accountBalanceDB.AccountId).Return(s.accountBalanceDB, sql.ErrNoRows)

	s.fnMock.On("FnSave", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Unexpected error"))
	s.fnMock.On("FnUpdate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := s.uc.Execute(s.ctx, s.inputDTO)

	s.NotNil(err)
	s.fnMock.AssertNumberOfCalls(s.T(), "FnSave", 1)
	s.fnMock.AssertNumberOfCalls(s.T(), "FnUpdate", 0)
	s.mockUow.AssertNumberOfCalls(s.T(), "GetRepository", 1)
	s.gatewayMock.AssertNumberOfCalls(s.T(), "Get", 1)
}

func (s *AcctBalTestSuite) TestReceiveAcctBalUseCase_ErrorFnUpdate() {
	s.mockUow.On("GetRepository", mock.Anything, mock.Anything).Return(s.gatewayMock, nil)
	s.gatewayMock.On("Get", s.accountBalanceDB.AccountId).Return(s.accountBalanceDB, nil)

	s.fnMock.On("FnUpdate", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Unexpected error"))
	s.fnMock.On("FnSave", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := s.uc.Execute(s.ctx, s.inputDTO)

	s.NotNil(err)
	s.fnMock.AssertNumberOfCalls(s.T(), "FnSave", 0)
	s.fnMock.AssertNumberOfCalls(s.T(), "FnUpdate", 1)
	s.mockUow.AssertNumberOfCalls(s.T(), "GetRepository", 1)
	s.gatewayMock.AssertNumberOfCalls(s.T(), "Get", 1)
}

func (s *AcctBalTestSuite) TestReceiveAcctBalUseCase_ErrorConvertInput() {
	s.inputDTO = CreateAccountBalanceInputDto{}
	err := s.uc.Execute(s.ctx, s.inputDTO)

	s.NotNil(err)
	s.Equal(err.Error(), "account id is required")

}

func (s *AcctBalTestSuite) TestReceiveAcctBalUseCase_ErrorUpdtEntityWithInputDTO() {
	s.inputDTO.DateTransaction = time.Time{}

	err := s.uc.Execute(s.ctx, s.inputDTO)

	s.NotNil(err)
	s.Equal(err.Error(), "date transaction is required")
}
