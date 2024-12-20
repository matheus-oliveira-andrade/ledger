package usecases

import (
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/slogger"
)

type GetAccountUseCaseInterface interface {
	Handle(accId string) (*domain.Account, error)
}

type GetAccountUseCaseImp struct {
	accountRepository repositories.AccountRepositoryInterface
	logger            slogger.LoggerInterface
}

func NewGetAccountUseCase(
	logger slogger.LoggerInterface,
	accountRepository repositories.AccountRepositoryInterface,
) GetAccountUseCaseInterface {
	return &GetAccountUseCaseImp{
		accountRepository: accountRepository,
		logger:            logger,
	}
}

func (us *GetAccountUseCaseImp) Handle(accId string) (*domain.Account, error) {
	us.logger.LogInformation("getting account by id", "accId", accId)

	acc, err := us.accountRepository.GetById(accId)

	if err != nil {
		us.logger.LogError("error getting account by id", "err", err)
		return nil, err
	}

	us.logger.LogInformation("searched account by id", "accId", accId)
	return acc, nil
}
