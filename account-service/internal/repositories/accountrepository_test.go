package repositories_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestCreate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewAccountRepository(db)

	acc := domain.NewAccount("acc test", "01234567890")

	mock.
		ExpectQuery("INSERT INTO accounts").
		WithArgs(acc.Name, acc.Document, acc.CreatedAt, acc.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"Id"}).AddRow("1"))

	id, err := repo.Create(acc)

	assert.NoError(t, err)
	assert.Equal(t, "1", id)
}

func TestCreate_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewAccountRepository(db)

	acc := domain.NewAccount("acc test", "01234567890")

	mock.
		ExpectQuery("INSERT INTO accounts").
		WithArgs(acc.Name, acc.Document, acc.CreatedAt, acc.UpdatedAt).
		WillReturnError(sqlmock.ErrCancelled)

	id, err := repo.Create(acc)

	assert.Error(t, err)
	assert.Empty(t, id)
}

func TestGetByDocument_AccFound(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	acc := domain.NewAccount("acc test", "01234567890")

	repo := repositories.NewAccountRepository(db)

	rows := sqlmock.
		NewRows([]string{"Id", "Name", "Document", "CreatedAt", "UpdatedAt"}).
		AddRow(acc.Id, acc.Name, acc.Document, acc.CreatedAt, acc.UpdatedAt)

	mock.
		ExpectQuery("SELECT Id, Name, Document, CreatedAt, UpdatedAt FROM accounts WHERE Document = \\$1").
		WithArgs(acc.Document).
		WillReturnRows(rows)

	account, err := repo.GetByDocument(acc.Document)

	assert.NoError(t, err)
	assert.NotNil(t, account)

	assert.Equal(t, acc, account)
}

func TestGetByDocument_AccNotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	defer db.Close()

	acc := domain.NewAccount("acc test", "01234567890")

	repo := repositories.NewAccountRepository(db)

	mock.
		ExpectQuery("SELECT Id, Name, Document, CreatedAt, UpdatedAt FROM accounts WHERE Document = \\$1").
		WithArgs(acc.Document).
		WillReturnError(sql.ErrNoRows)

	account, err := repo.GetByDocument(acc.Document)

	assert.NoError(t, err)
	assert.Nil(t, account)
}
