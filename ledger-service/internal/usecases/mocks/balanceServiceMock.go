package usecases_mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockBalanceService struct {
	mock.Mock
}

func (m *MockBalanceService) CalculateBalance(ctx context.Context, accountID int64) (int64, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).(int64), args.Error(1)
}
