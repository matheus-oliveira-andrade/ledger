package usecases

import (
	"errors"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils/slogger"

	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/repositories"
)

type CreateAccountUseCaseInterface interface {
	Handle(document, name string) (string, error)
}

type CreateAccountUseCaseImp struct {
	logger            slogger.LoggerInterface
	accountRepository repositories.AccountRepositoryInterface
}

func NewCreateAccountUseCase(
	logger slogger.LoggerInterface,
	accountRepository repositories.AccountRepositoryInterface) *CreateAccountUseCaseImp {
	return &CreateAccountUseCaseImp{
		logger:            logger,
		accountRepository: accountRepository,
	}
}

func (u *CreateAccountUseCaseImp) Handle(document, name string) (string, error) {
	u.logger.LogInformation("creating account", "document", document)

	acc, err := u.accountRepository.GetByDocument(document)
	if err != nil {
		u.logger.LogError("error getting account", "error", err.Error())
		return "", err
	}

	if acc != nil {
		u.logger.LogError("account already exists", "document", document)
		return "", errors.New("account already exists")
	}

	acc = domain.NewAccount(name, document)

	accId, err := u.accountRepository.Create(acc)
	if err != nil {
		u.logger.LogError("error creating account", "error", err.Error())
		return "", err
	}

	u.logger.LogInformation("account created", "accId", accId, "document", document)
	return accId, nil
}
