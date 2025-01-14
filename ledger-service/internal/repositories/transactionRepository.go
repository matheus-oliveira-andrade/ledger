package repositories

import (
	"database/sql"
	"strconv"

	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
)

type TransactionRepositoryInterface interface {
	Create(acc *domain.Transaction) (int64, error)
}

type TransactionRepositoryImp struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepositoryImp {
	return &TransactionRepositoryImp{
		db: db,
	}
}

func (r *TransactionRepositoryImp) Create(transaction *domain.Transaction) (int64, error) {
	row := r.db.QueryRow(`
	INSERT INTO Transaction (Description, CreatedAt)
	VALUES ($1, $2)
	
	RETURNING Id
	`, transaction.Description, transaction.CreatedAt)

	var id string
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, err
	}

	return idInt, err
}
