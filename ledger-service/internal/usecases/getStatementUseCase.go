package usecases

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/slogger"
	"strings"
)

type GetStatementUseCaseInterface interface {
	Handle(accId int64, transactionsType string, period int, limit int, page int) (*domain.Statement, error)
}

type GetStatementUseCaseImp struct {
	logger              slogger.LoggerInterface
	statementRepository repositories.StatementRepositoryInterface
}

func NewGetStatementUseCase(
	logger slogger.LoggerInterface,
	statementRepository repositories.StatementRepositoryInterface) *GetStatementUseCaseImp {
	return &GetStatementUseCaseImp{
		logger:              logger,
		statementRepository: statementRepository,
	}
}

func (us *GetStatementUseCaseImp) Handle(accId int64, transactionsType string, period int, limit int, page int) (*domain.Statement, error) {
	us.logger.LogInformation("getting statements", "accId", accId, "transactionsType", transactionsType, "period", period)

	startDate, endDate := domain.GetStatementPeriodDates(period)
	entriesType := us.getEntriesType(transactionsType)

	transactions, hasNextPage, err := us.statementRepository.GetStatementTransactions(accId, startDate, endDate, entriesType, limit, page)
	if err != nil {
		us.logger.LogError("error getting statement transactions", err)
		return nil, err
	}

	statement := domain.NewStatement(accId, startDate, endDate, transactions, hasNextPage)
	return statement, nil
}

func (us *GetStatementUseCaseImp) getEntriesType(entriesType string) []string {
	if entriesType == "" || entriesType == "ALL" {
		return []string{"DEBIT", "CREDIT"}
	}

	return strings.Split(entriesType, ",")
}
