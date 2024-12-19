package usecases_mocks

import (
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(acc *domain.Account) (string, error) {
	args := m.Called(acc)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockAccountRepository) GetByDocument(document string) (*domain.Account, error) {
	args := m.Called(document)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetById(id string) (*domain.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Account), args.Error(1)
}
