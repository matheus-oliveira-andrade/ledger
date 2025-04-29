package domain

import "time"

type StatementTransaction struct {
	AccountId   int64
	Description string
	Amount      int64
	EntryType   EntryType
	CreatedAt   time.Time
}
