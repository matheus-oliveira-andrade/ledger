package domain_test

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"testing"
)

func TestBalanceCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name          string
		lines         []domain.TransactionLine
		expected      int64
		mustHaveError bool
	}{
		{
			name: "should calculate the balance correctly",
			lines: []domain.TransactionLine{
				{EntryType: domain.Credit, Amount: 100},
				{EntryType: domain.Debit, Amount: 50},
			},
			expected:      50,
			mustHaveError: false,
		},
		{
			name: "should return an error for unmapped entry type",
			lines: []domain.TransactionLine{
				{EntryType: "unmapped", Amount: 100},
			},
			expected:      0,
			mustHaveError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			bc := domain.NewBalanceCalculator()
			result, err := bc.Calculate(test.lines)

			if test.mustHaveError && err == nil {
				t.Errorf("expected an error but got nil")
			}

			if result != test.expected {
				t.Errorf("expected %d but got %d", test.expected, result)
			}
		})
	}
}
