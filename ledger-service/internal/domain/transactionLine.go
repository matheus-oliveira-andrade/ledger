package domain

import "time"

const (
	Debit  EntryType = "DEBIT"
	Credit EntryType = "CREDIT"
)

type EntryType string

type TransactionLine struct {
	Id            int64
	AccountId     int64
	TransactionId int64
	Amount        int64
	EntryType     EntryType
	CreatedAt     time.Time
}

func NewTransactionLine(accountId, amount int64, entryType EntryType) *TransactionLine {
	return &TransactionLine{
		AccountId: accountId,
		Amount:    amount,
		EntryType: entryType,
	}
}
