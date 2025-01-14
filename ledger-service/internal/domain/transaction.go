package domain

import "time"

type Transaction struct {
	Id          int64
	Description string
	CreatedAt   time.Time
	Lines       []*TransactionLine
}

func NewTransaction(description string, lines []*TransactionLine) *Transaction {
	return &Transaction{
		Description: description,
		CreatedAt:   time.Now(),
		Lines:       lines,
	}
}
