package usecases

import (
	"context"
	"errors"
	accountgrpc "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/adapters/grpc"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/utils/slogger"
	"strconv"

	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/services"
)

type FundsTransferUseCaseInterface interface {
	Handle(ctx context.Context, accFrom, accTo, amount int64) error
}

var (
	ErrFromAccountNotFound = errors.New("acc from not found")
	ErrToAccountNotFound   = errors.New("acc to not found")
	ErrCreatingTransaction = errors.New("error persisting transaction")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

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
	u.logger.LogInformationContext(ctx, "processing transfer of funds", "accFrom", accFrom, "accTo", accTo, "amount", amount)

	err := u.validateFromAccount(ctx, accFrom, amount)
	if err != nil {
		return err
	}

	err = u.validateToAccount(ctx, accTo)
	if err != nil {
		return err
	}

	err = u.createTransaction(ctx, accFrom, accTo, amount)
	if err != nil {
		return err
	}

	u.logger.LogInformationContext(ctx, "transfer completed")
	return nil
}

func (u *FundsTransferUseCaseImp) accountExist(ctx context.Context, accId int64) (bool, error) {
	req := accountgrpc.GetAccountRequest{
		AccId: strconv.FormatInt(accId, 10),
	}

	u.logger.LogInformationContext(ctx, "searching for account in account server", "accId", req.AccId)

	acc, err := u.accountClient.GetAccount(ctx, &req)

	u.logger.LogInformationContext(ctx, "searched for account in account server", "accId", req.AccId, "acc", acc, "err", err)

	return acc != nil && acc.Id == req.AccId, err
}

func (u *FundsTransferUseCaseImp) getBalance(ctx context.Context, accId int64) (int64, error) {
	return u.balanceService.CalculateBalance(ctx, accId)
}

func (u *FundsTransferUseCaseImp) validateFromAccount(ctx context.Context, accFrom int64, amount int64) error {
	ok, err := u.accountExist(ctx, accFrom)
	if err != nil {
		u.logger.LogErrorContext(ctx, "error checking if account exist", "accFrom", accFrom, "error", err)
		return err
	}

	if !ok {
		u.logger.LogErrorContext(ctx, "acc from not found", "accFrom", accFrom)
		return ErrFromAccountNotFound
	}

	balance, err := u.getBalance(ctx, accFrom)
	if err != nil {
		u.logger.LogErrorContext(ctx, "error get acc balance", "accFrom", accFrom, "error", err)
		return err
	}

	if balance < amount {
		u.logger.LogErrorContext(ctx, "insufficient balance", "accFrom", accFrom)
		return ErrInsufficientBalance
	}

	return nil
}

func (u *FundsTransferUseCaseImp) validateToAccount(ctx context.Context, accTo int64) error {
	ok, err := u.accountExist(ctx, accTo)
	if err != nil {
		u.logger.LogErrorContext(ctx, "error checking if account exist", "accTo", accTo, "error", err)
		return err
	}

	if !ok {
		u.logger.LogErrorContext(ctx, "acc to not found", "accTo", accTo)
		return ErrToAccountNotFound
	}

	return nil
}

func (u *FundsTransferUseCaseImp) createTransaction(ctx context.Context, accFrom int64, accTo int64, amount int64) error {
	lineFrom := domain.NewTransactionLine(accFrom, amount, domain.Debit)
	lineTo := domain.NewTransactionLine(accTo, amount, domain.Credit)

	transaction := domain.NewTransaction("transfer between accounts", []*domain.TransactionLine{lineFrom, lineTo})

	err := u.transactionService.Save(transaction)
	if err != nil {
		u.logger.LogErrorContext(ctx, "error persisting transaction", "accFrom", accFrom, "accTo", accTo, "error", err)
		return ErrCreatingTransaction
	}

	return nil
}
