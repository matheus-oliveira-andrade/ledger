package usecases_mocks

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) Save(transaction *domain.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}
