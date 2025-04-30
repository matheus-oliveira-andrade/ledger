package usecases_mocks

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockStatementRepository struct {
	mock.Mock
}

func (s *MockStatementRepository) GetStatementTransactions(accId int64, startDate, endDate time.Time, entriesType []string, limit, page int) (*[]domain.StatementTransaction, bool, error) {
	args := s.Called(accId, startDate, endDate, entriesType, limit, page)
	return args.Get(0).(*[]domain.StatementTransaction), args.Bool(1), args.Error(2)
}
