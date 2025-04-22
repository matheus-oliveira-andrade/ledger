package services_mocks

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockTransactionLineRepository struct {
	mock.Mock
}

func (m *MockTransactionLineRepository) Create(line *domain.TransactionLine) (string, error) {
	args := m.Called(line)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockTransactionLineRepository) GetTransactions(accId int64) (*[]domain.TransactionLine, error) {
	args := m.Called(accId)
	return args.Get(0).(*[]domain.TransactionLine), args.Error(1)
}
