package usecases_test

import (
	"errors"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/usecases"
	usecases_mocks "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestHandle_Error(t *testing.T) {
	t.Parallel()

	// Arrange
	accountId := int64(12345)

	transactionsType := "ALL"
	entriesType := []string{"DEBIT", "CREDIT"}

	mockStatementRepository := usecases_mocks.MockStatementRepository{}
	mockStatementRepository.
		On("GetStatementTransactions", accountId, mock.Anything, mock.Anything, entriesType, 2, 0).
		Return(&[]domain.StatementTransaction{}, false, errors.New("error"))

	mockLogger := usecases_mocks.MockLogger{}
	mockLogger.
		On("LogInformation", mock.Anything, mock.Anything).
		Return()
	mockLogger.
		On("LogError", mock.Anything, mock.Anything).
		Return()

	useCase := usecases.NewGetStatementUseCase(&mockLogger, &mockStatementRepository)

	// Act
	statement, err := useCase.Handle(accountId, transactionsType, 60, 2, 0)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, statement)

	mockStatementRepository.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestHandle_NotFoundTransactions(t *testing.T) {
	t.Parallel()

	// Arrange
	accountId := int64(12345)

	transactionsType := "ALL"
	entriesType := []string{"DEBIT", "CREDIT"}

	mockStatementRepository := usecases_mocks.MockStatementRepository{}
	mockStatementRepository.
		On("GetStatementTransactions", accountId, mock.Anything, mock.Anything, entriesType, 2, 0).
		Return(&[]domain.StatementTransaction{}, false, nil)

	mockLogger := usecases_mocks.MockLogger{}
	mockLogger.
		On("LogInformation", mock.Anything, mock.Anything).
		Return()

	useCase := usecases.NewGetStatementUseCase(&mockLogger, &mockStatementRepository)

	// Act
	statement, err := useCase.Handle(accountId, transactionsType, 60, 2, 0)

	// Assert
	assert.Nil(t, err)

	assert.NotNil(t, statement)
	assert.Empty(t, statement.Transactions)
	assert.Equal(t, accountId, statement.AccountId)
	assert.False(t, statement.HasNextPage)

	mockStatementRepository.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestHandle_FoundTransactions(t *testing.T) {
	t.Parallel()

	// Arrange
	accountId := int64(12345)

	transactionsType := "ALL"
	entriesType := []string{"DEBIT", "CREDIT"}

	transactions := &[]domain.StatementTransaction{
		{
			AccountId:   accountId,
			Amount:      100.0,
			EntryType:   domain.Debit,
			CreatedAt:   time.Now().AddDate(0, 0, -10),
			Description: "Transaction 1",
		},
	}

	mockStatementRepository := usecases_mocks.MockStatementRepository{}
	mockStatementRepository.
		On("GetStatementTransactions", accountId, mock.Anything, mock.Anything, entriesType, 2, 0).
		Return(transactions, false, nil)

	mockLogger := usecases_mocks.MockLogger{}
	mockLogger.
		On("LogInformation", mock.Anything, mock.Anything).
		Return()

	useCase := usecases.NewGetStatementUseCase(&mockLogger, &mockStatementRepository)

	// Act
	statement, err := useCase.Handle(accountId, transactionsType, 60, 2, 0)

	// Assert
	assert.Nil(t, err)

	assert.NotNil(t, statement)
	assert.Len(t, *statement.Transactions, 1)
	assert.Equal(t, accountId, statement.AccountId)
	assert.False(t, statement.HasNextPage)

	mockStatementRepository.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestHandle_FoundTransactions_MustHaveNextPage(t *testing.T) {
	t.Parallel()

	// Arrange
	accountId := int64(12345)

	transactionsType := "ALL"
	entriesType := []string{"DEBIT", "CREDIT"}

	transactions := &[]domain.StatementTransaction{
		{
			AccountId:   accountId,
			Amount:      100.0,
			EntryType:   domain.Debit,
			CreatedAt:   time.Now().AddDate(0, 0, -10),
			Description: "Transaction 1",
		},
		{
			AccountId:   accountId,
			Amount:      50.0,
			EntryType:   domain.Credit,
			CreatedAt:   time.Now().AddDate(0, 0, -11),
			Description: "Transaction 2",
		},
	}

	mockStatementRepository := usecases_mocks.MockStatementRepository{}
	mockStatementRepository.
		On("GetStatementTransactions", accountId, mock.Anything, mock.Anything, entriesType, 2, 0).
		Return(transactions, true, nil)

	mockLogger := usecases_mocks.MockLogger{}
	mockLogger.
		On("LogInformation", mock.Anything, mock.Anything).
		Return()

	useCase := usecases.NewGetStatementUseCase(&mockLogger, &mockStatementRepository)

	// Act
	statement, err := useCase.Handle(accountId, transactionsType, 60, 2, 0)

	// Assert
	assert.Nil(t, err)

	assert.NotNil(t, statement)
	assert.Len(t, *statement.Transactions, 2)
	assert.Equal(t, accountId, statement.AccountId)
	assert.True(t, statement.HasNextPage)

	mockStatementRepository.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
