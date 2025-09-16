package domain

import (
	"time"
)

type Transaction struct {
	Id          int64
	Description string
	CreatedAt   time.Time
	lines       []*TransactionLine
}

func NewTransaction(amount int64, description string, accFrom, accTo *Account) *Transaction {
	lineFrom := NewTransactionLine(accFrom.Id, amount, Debit)
	lineTo := NewTransactionLine(accTo.Id, amount, Credit)

	return &Transaction{
		Description: description,
		CreatedAt:   time.Now(),
		lines: []*TransactionLine{
			lineFrom,
			lineTo,
		},
	}
}

func (t *Transaction) GetLines() []*TransactionLine {
	return t.lines
}
