package domain

import "time"

func GetStatementPeriodDates(period int) (time.Time, time.Time) {
	startDate := time.Now().AddDate(0, 0, period*-1)
	endDate := time.Now().AddDate(0, 0, 1)

	return truncateTime(startDate), truncateTime(endDate)
}

func truncateTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
