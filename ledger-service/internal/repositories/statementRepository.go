package repositories

import (
	"database/sql"
	"github.com/lib/pq"
	"time"

	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
)

type StatementRepositoryInterface interface {
	GetStatementTransactions(accId int64, startDate time.Time, endDate time.Time, entriesType []string, limit int, page int) (*[]domain.StatementTransaction, bool, error)
}

type StatementRepositoryImp struct {
	db *sql.DB
}

func NewStatementRepository(db *sql.DB) *StatementRepositoryImp {
	return &StatementRepositoryImp{
		db: db,
	}
}
func (r *StatementRepositoryImp) GetStatementTransactions(accId int64, startDate time.Time, endDate time.Time, entriesType []string, limit int, page int) (*[]domain.StatementTransaction, bool, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := r.calculateOffset(page, limit)

	rows, err := r.db.Query(`
		SELECT tl.AccountId, t.Description, tl.Amount, tl.EntryType, t.CreatedAt
		FROM TransactionLine tl
		INNER JOIN Transaction t ON tl.TransactionId = t.Id		
		WHERE tl.AccountId = $1 
	      AND t.CreatedAt BETWEEN $2 AND $3 
		  AND tl.EntryType = ANY($4)
		ORDER BY t.CreatedAt DESC
		LIMIT $5 OFFSET $6
	`, accId, startDate, endDate, pq.Array(entriesType), limit+1, offset)

	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var result []domain.StatementTransaction

	for rows.Next() {
		var line domain.StatementTransaction

		err := rows.Scan(&line.AccountId, &line.Description, &line.Amount, &line.EntryType, &line.CreatedAt)
		if err != nil {
			return nil, false, err
		}

		result = append(result, line)
	}

	if err = rows.Err(); err != nil {
		return nil, false, err
	}

	if len(result) > limit {
		result = result[:limit]
		return &result, true, nil
	}

	return &result, false, nil
}

func (r *StatementRepositoryImp) calculateOffset(page int, limit int) int {
	return (page - 1) * limit
}
