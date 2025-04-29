package repositories_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetStatementTransactions_Error(t *testing.T) {
	t.Parallel()

	// Arrange
	db, mock, err := sqlmock.New()
	defer db.Close()

	accId := int64(12345)
	startDate := time.Now()
	endDate := time.Now().AddDate(0, 0, 1)
	entriesType := "DEBIT,CREDIT"

	mock.
		ExpectQuery("SELECT tl.AccountId, t.Description, tl.Amount, tl.EntryType, t.CreatedAt FROM TransactionLine tl INNER JOIN Transaction t ON tl.TransactionId = t.Id WHERE tl.AccountId = \\$1 AND t.CreatedAt BETWEEN \\$2 AND \\$3 AND tl.EntryType IN (\\$4)").
		WithArgs(accId, startDate, endDate, entriesType, 2, 1).
		WillReturnError(sql.ErrConnDone)

	repository := repositories.NewStatementRepository(db)

	// Act
	result, hasNextPage, err := repository.GetStatementTransactions(accId, startDate, endDate, []string{"DEBIT", "CREDIT"}, 2, 0)

	// Assert
	assert.Error(t, sql.ErrConnDone, err)
	assert.Nil(t, result)
	assert.False(t, hasNextPage)
}

func TestGetStatementTransactions_NoRows(t *testing.T) {
	t.Parallel()

	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	accId := int64(12345)
	startDate := time.Now()
	endDate := time.Now().AddDate(0, 0, 2)
	entriesType := "DEBIT,CREDIT"

	mock.
		ExpectQuery("SELECT tl.AccountId, t.Description, tl.Amount, tl.EntryType, t.CreatedAt").
		WithArgs(accId, startDate, endDate, entriesType, 3, 1).
		WillReturnRows(sqlmock.NewRows([]string{"AccountId", "Description", "Amount", "EntryType", "CreatedAt"}))

	repository := repositories.NewStatementRepository(db)

	// Act
	result, hasNextPage, err := repository.GetStatementTransactions(accId, startDate, endDate, []string{"DEBIT", "CREDIT"}, 2, 0)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.False(t, hasNextPage)
}

func TestGetStatementTransactions_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	accId := int64(12345)
	startDate := time.Now()
	endDate := time.Now().AddDate(0, 0, 2)
	entriesType := "DEBIT,CREDIT"

	rows := sqlmock.
		NewRows([]string{"AccountId", "Description", "Amount", "EntryType", "CreatedAt"}).
		AddRow(accId, "Test Transaction", 100, "DEBIT", startDate).
		AddRow(accId, "Test Transaction", 100, "DEBIT", startDate)

	mock.
		ExpectQuery("SELECT tl.AccountId, t.Description, tl.Amount, tl.EntryType, t.CreatedAt").
		WithArgs(accId, startDate, endDate, entriesType, 2, 0).
		WillReturnRows(rows)

	repository := repositories.NewStatementRepository(db)

	// Act
	result, hasNextPage, err := repository.GetStatementTransactions(accId, startDate, endDate, []string{"DEBIT", "CREDIT"}, 1, 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, len(*result))
	assert.True(t, hasNextPage)
}
