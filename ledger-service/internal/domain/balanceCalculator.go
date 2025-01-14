package domain

import (
	"fmt"
)

type BalanceCalculator struct {
}

func (bc *BalanceCalculator) Calculate(lines []TransactionLine) (int64, error) {
	var creditsSum int64 = 0
	var debitsSum int64 = 0

	for _, line := range lines {

		if line.EntryType == Credit {
			creditsSum += line.Amount
			continue
		}

		if line.EntryType == Debit {
			debitsSum += line.Amount
			continue
		}

		return 0, fmt.Errorf("unmapped entry type (%v)", line.EntryType)
	}

	return creditsSum - debitsSum, nil
}
