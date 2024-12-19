package repositories

import (
	"database/sql"

	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/domain"
)

type AccountRepositoryInterface interface {
	Create(acc *domain.Account) (string, error)
	GetByDocument(document string) (*domain.Account, error)
	GetById(id string) (*domain.Account, error)
}

type AccountRepositoryImp struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepositoryImp {
	return &AccountRepositoryImp{
		db: db,
	}
}

func (r *AccountRepositoryImp) Create(acc *domain.Account) (string, error) {
	row := r.db.QueryRow(`
	INSERT INTO accounts (Name, Document, CreatedAt, UpdatedAt)
	VALUES ($1, $2, $3, $4)
	
	RETURNING Id
	`, acc.Name, acc.Document, acc.CreatedAt, acc.UpdatedAt)

	var id string
	err := row.Scan(&id)

	return id, err
}

func (r *AccountRepositoryImp) GetByDocument(document string) (*domain.Account, error) {
	row := r.db.QueryRow(`
		SELECT Id, Name, Document, CreatedAt, UpdatedAt
		FROM accounts 
		WHERE Document = $1
	`, document)

	var account domain.Account
	err := row.Scan(&account.Id, &account.Name, &account.Document, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}

func (r *AccountRepositoryImp) GetById(id string) (*domain.Account, error) {
	row := r.db.QueryRow(`
		SELECT Id, Name, Document, CreatedAt, UpdatedAt
		FROM accounts 
		WHERE Id = $1
	`, id)

	var account domain.Account
	err := row.Scan(&account.Id, &account.Name, &account.Document, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &account, nil
}
