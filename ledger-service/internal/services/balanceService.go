package services

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/repositories"
)

type BalanceServiceInterface interface {
	CalculateBalance(accId int64) (int64, error)
}

type BalanceServiceImp struct {
	transactionLineRepository repositories.TransactionLineRepositoryInterface
}

func NewBalanceService(transactionLineRepository repositories.TransactionLineRepositoryInterface) *BalanceServiceImp {
	return &BalanceServiceImp{
		transactionLineRepository: transactionLineRepository,
	}
}

func (s *BalanceServiceImp) CalculateBalance(accId int64) (int64, error) {
	lines, err := s.transactionLineRepository.GetTransactions(accId)
	if err != nil {
		return 0, err
	}

	return domain.NewBalanceCalculator().Calculate(*lines)
}
