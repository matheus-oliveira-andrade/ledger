package adapters_mocks

import (
	"context"

	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountAdapter struct {
	mock.Mock
}

func (m *MockAccountAdapter) GetAccount(ctx context.Context, accId int64) (*domain.Account, error) {
	args := m.Called(ctx, accId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}
