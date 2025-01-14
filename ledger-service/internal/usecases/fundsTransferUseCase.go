package usecases

import (
	"context"
	"errors"
	"strconv"

	accountgrpc "github.com/matheus-oliveira-andrade/ledger/ledger-service/grpc"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/services"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/slogger"
)

type FundsTransferUseCaseInterface interface {
	Handle(ctx context.Context, accFrom, accTo, amount int64) error
}

type FundsTransferUseCaseImp struct {
	logger             slogger.LoggerInterface
	transactionService services.TransactionServiceInterface
	accountClient      accountgrpc.AccountClient
}

func NewFundsTransferUseCase(
	logger slogger.LoggerInterface,
	transactionService services.TransactionServiceInterface,
	accountClient accountgrpc.AccountClient) *FundsTransferUseCaseImp {
	return &FundsTransferUseCaseImp{
		logger:             logger,
		transactionService: transactionService,
		accountClient:      accountClient,
	}
}

func (u *FundsTransferUseCaseImp) Handle(ctx context.Context, accFrom, accTo, amount int64) error {
	u.logger.LogInformation("processing transfer of funds", "accFrom", accFrom, "accTo", accTo, "amount", amount)

	ok, err := u.accountExist(ctx, accFrom)
	if err != nil {
		u.logger.LogError("error checking if account exist", "accFrom", accFrom, "error", err)
		return err
	}
	if !ok {
		u.logger.LogError("acc from not found", "accFrom", accFrom)
		return errors.New("acc from not found")
	}

	ok, err = u.accountExist(ctx, accTo)
	if err != nil {
		u.logger.LogError("error checking if account exist", "accTo", accTo, "error", err)
		return err
	}
	if !ok {
		u.logger.LogError("acc to not found", "accTo", accTo)
		return errors.New("acc to not found")
	}

	lineFrom := domain.NewTransactionLine(accFrom, amount, domain.Debit)
	lineTo := domain.NewTransactionLine(accTo, amount, domain.Credit)

	transaction := domain.NewTransaction("transfer between accounts", []*domain.TransactionLine{lineFrom, lineTo})

	u.transactionService.Save(transaction)
	u.logger.LogInformation("transfer completed")

	return nil
}

func (u *FundsTransferUseCaseImp) accountExist(ctx context.Context, accId int64) (bool, error) {
	req := accountgrpc.GetAccountRequest{
		AccId: strconv.FormatInt(int64(accId), 10),
	}

	u.logger.LogError("searching for acccount in account server", "accId", req.AccId)

	acc, err := u.accountClient.GetAccount(ctx, &req)

	u.logger.LogError("searched for acccount in account server", "accId", req.AccId, "acc", acc, "err", err)

	return acc != nil && acc.Id == req.AccId, err
}
