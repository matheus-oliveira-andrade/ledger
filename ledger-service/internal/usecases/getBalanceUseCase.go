package usecases

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/services"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/slogger"
)

type GetBalanceUseCaseInterface interface {
	Handle(accId int64) (int64, error)
}

type GetBalanceUseCaseImp struct {
	logger         slogger.LoggerInterface
	balanceService services.BalanceServiceInterface
}

func NewGetBalanceUseCase(
	logger slogger.LoggerInterface,
	balanceService services.BalanceServiceInterface) *GetBalanceUseCaseImp {
	return &GetBalanceUseCaseImp{
		logger:         logger,
		balanceService: balanceService,
	}
}

func (us *GetBalanceUseCaseImp) Handle(accId int64) (int64, error) {
	return us.balanceService.CalculateBalance(accId)
}
