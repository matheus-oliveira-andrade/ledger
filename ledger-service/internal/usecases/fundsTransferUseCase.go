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
	balanceService     services.BalanceServiceInterface
}

func NewFundsTransferUseCase(
	logger slogger.LoggerInterface,
	transactionService services.TransactionServiceInterface,
	accountClient accountgrpc.AccountClient,
	balanceService services.BalanceServiceInterface) *FundsTransferUseCaseImp {
	return &FundsTransferUseCaseImp{
		logger:             logger,
		transactionService: transactionService,
		accountClient:      accountClient,
		balanceService:     balanceService,
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

	balance, err := u.getBalance(accFrom)
	if err != nil {
		u.logger.LogError("error get acc balance", "accFrom", accFrom, "error", err)
		return err
	}
	if balance < amount {
		u.logger.LogError("insufficient balance", "accFrom", accFrom)
		return errors.New("insufficient balance")
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

	err = u.transactionService.Save(transaction)
	if err != nil {
		u.logger.LogError("error persisting transaction", "accFrom", accFrom, "accTo", accTo, "error", err)
		return errors.New("error persisting transaction")
	}

	u.logger.LogInformation("transfer completed")
	return nil
}

func (u *FundsTransferUseCaseImp) accountExist(ctx context.Context, accId int64) (bool, error) {
	req := accountgrpc.GetAccountRequest{
		AccId: strconv.FormatInt(accId, 10),
	}

	u.logger.LogError("searching for account in account server", "accId", req.AccId)

	acc, err := u.accountClient.GetAccount(ctx, &req)

	u.logger.LogError("searched for account in account server", "accId", req.AccId, "acc", acc, "err", err)

	return acc != nil && acc.Id == req.AccId, err
}

func (u *FundsTransferUseCaseImp) getBalance(accId int64) (int64, error) {
	return u.balanceService.CalculateBalance(accId)
}
