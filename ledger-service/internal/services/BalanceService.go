package services

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/repositories"
)

type BalanceServiceInterface interface {
	CalculateBalance(accId int64) (int64, error)
}

type BalanceService struct {
	transactionLineRepository repositories.TransactionLineRepositoryInterface
}

func (s *BalanceService) CalculateBalance(accId int64) (int64, error) {
	lines, err := s.transactionLineRepository.GetTransactions(accId)
	if err != nil {
		return 0, err
	}

	balanceCalculator := domain.BalanceCalculator{}
	return balanceCalculator.Calculate(*lines)
}
