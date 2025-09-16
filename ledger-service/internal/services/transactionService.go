package services

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/utils/slogger"
)

type TransactionServiceInterface interface {
	Save(transaction *domain.Transaction) error
}

type TransactionService struct {
	logger                    slogger.LoggerInterface
	transactionRepository     repositories.TransactionRepositoryInterface
	transactionLineRepository repositories.TransactionLineRepositoryInterface
}

func NewTransactionService(
	logger slogger.LoggerInterface,
	transactionRepository repositories.TransactionRepositoryInterface,
	transactionLineRepository repositories.TransactionLineRepositoryInterface) *TransactionService {
	return &TransactionService{
		logger:                    logger,
		transactionRepository:     transactionRepository,
		transactionLineRepository: transactionLineRepository,
	}
}

func (ts *TransactionService) Save(transaction *domain.Transaction) error {
	ts.logger.LogInformation("saving transaction")

	transactionId, err := ts.transactionRepository.Create(transaction)
	if err != nil {
		return err
	}

	for _, line := range transaction.GetLines() {
		line.TransactionId = transactionId

		_, err := ts.transactionLineRepository.Create(line)
		if err != nil {
			return err
		}
	}

	return nil
}
