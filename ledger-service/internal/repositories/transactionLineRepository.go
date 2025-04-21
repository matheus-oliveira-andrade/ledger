package repositories

import (
	"database/sql"

	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
)

type TransactionLineRepositoryInterface interface {
	Create(line *domain.TransactionLine) (string, error)
	GetTransactions(accId int64) (*[]domain.TransactionLine, error)
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

func (r *TransactionLineRepositoryImp) GetTransactions(accId int64) (*[]domain.TransactionLine, error) {
	rows, err := r.db.Query(`
		SELECT Id, AccountId, TransactionId, Amount, EntryType, CreatedAt
		FROM TransactionLine
		WHERE AccountId = $1
	`, accId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []domain.TransactionLine
	for rows.Next() {
		var line domain.TransactionLine
		err := rows.Scan(&line.Id, &line.AccountId, &line.TransactionId, &line.Amount, &line.EntryType, &line.CreatedAt)

		if err != nil {
			return nil, err
		}

		lines = append(lines, line)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &lines, nil
}
