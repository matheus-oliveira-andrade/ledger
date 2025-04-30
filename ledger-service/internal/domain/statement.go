package domain

import "time"

type Statement struct {
	AccountId    int64
	StartDate    time.Time
	EndDate      time.Time
	Transactions *[]StatementTransaction
	HasNextPage  bool
}

func NewStatement(accId int64, startDate time.Time, endDate time.Time, transactions *[]StatementTransaction, hasNextPage bool) *Statement {
	return &Statement{
		AccountId:    accId,
		StartDate:    startDate,
		EndDate:      endDate,
		Transactions: transactions,
		HasNextPage:  hasNextPage,
	}
}
