package usecases_mocks

import "github.com/stretchr/testify/mock"

type MockBalanceService struct {
	mock.Mock
}

func (m *MockBalanceService) CalculateBalance(accountID int64) (int64, error) {
	args := m.Called(accountID)
	return args.Get(0).(int64), args.Error(1)
}
