package domain

import "time"

func GetStatementPeriodDates(period int) (time.Time, time.Time) {
	startDate := time.Now().AddDate(0, 0, period*-1)
	endDate := time.Now()

	return startDate, endDate
}
