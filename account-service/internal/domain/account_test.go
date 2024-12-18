package domain_test

import (
	"testing"

	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName string
		name     string
		number   string
		document string
	}{
		{
			"should create a new account",
			"name test",
			"123",
			"01234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			acc := domain.NewAccount(tt.name, tt.document)

			assert.Equal(t, tt.name, acc.Name)
			assert.Equal(t, tt.document, acc.Document)
		})
	}

}
