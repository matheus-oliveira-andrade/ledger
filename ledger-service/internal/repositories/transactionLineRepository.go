package repositories

import (
	"database/sql"

	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
)

type TransactionLineRepositoryInterface interface {
	Create(acc *domain.TransactionLine) (string, error)
}

type TransactionLineRepositoryImp struct {
	db *sql.DB
}

func NewTransactionLineRepository(db *sql.DB) *TransactionLineRepositoryImp {
	return &TransactionLineRepositoryImp{
		db: db,
	}
}

func (r *TransactionLineRepositoryImp) Create(line *domain.TransactionLine) (string, error) {
	row := r.db.QueryRow(`
	INSERT INTO TransactionLine (AccountId, TransactionId, Amount, EntryType, CreatedAt)
	VALUES ($1, $2, $3, $4, $5)
	
	RETURNING Id
	`, line.AccountId, line.TransactionId, line.Amount, line.EntryType, line.CreatedAt)

	var id string
	err := row.Scan(&id)

	return id, err
}
