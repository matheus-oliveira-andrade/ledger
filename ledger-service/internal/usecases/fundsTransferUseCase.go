package usecases

import (
	"context"
	"errors"

	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/adapters"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/utils/slogger"

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
	accountAdapter     adapters.AccountAdapterInterface
	balanceService     services.BalanceServiceInterface
}

func NewFundsTransferUseCase(
	logger slogger.LoggerInterface,
	transactionService services.TransactionServiceInterface,
	accountAdapter adapters.AccountAdapterInterface,
	balanceService services.BalanceServiceInterface) *FundsTransferUseCaseImp {
	return &FundsTransferUseCaseImp{
		logger:             logger,
		transactionService: transactionService,
		accountAdapter:     accountAdapter,
		balanceService:     balanceService,
	}
}

func (u *FundsTransferUseCaseImp) Handle(ctx context.Context, accIdFrom, accIdTo, amount int64) error {
	u.logger.LogInformationContext(ctx, "processing transfer of funds", "accIdFrom", accIdFrom, "accIdTo", accIdTo, "amount", amount)

	accFrom, err := u.validateFromAccount(ctx, accIdFrom, amount)
	if err != nil {
		return err
	}

	accTo, err := u.validateToAccount(ctx, accIdTo)
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

func (u *FundsTransferUseCaseImp) getBalance(ctx context.Context, accId int64) (int64, error) {
	return u.balanceService.CalculateBalance(ctx, accId)
}

func (u *FundsTransferUseCaseImp) validateFromAccount(ctx context.Context, accIdFrom int64, amount int64) (*domain.Account, error) {
	account, err := u.accountAdapter.GetAccount(ctx, accIdFrom)

	if err != nil {
		u.logger.LogErrorContext(ctx, "error checking if account exist", "accIdFrom", accIdFrom, "error", err)
		return nil, err
	}

	if account == nil {
		u.logger.LogErrorContext(ctx, "acc from not found", "accIdFrom", accIdFrom)
		return nil, ErrFromAccountNotFound
	}

	balance, err := u.getBalance(ctx, accIdFrom)
	if err != nil {
		u.logger.LogErrorContext(ctx, "error get acc balance", "accIdFrom", accIdFrom, "error", err)
		return nil, err
	}

	if balance < amount {
		u.logger.LogErrorContext(ctx, "insufficient balance", "accIdFrom", accIdFrom)
		return nil, ErrInsufficientBalance
	}

	return account, nil
}

func (u *FundsTransferUseCaseImp) validateToAccount(ctx context.Context, accIdTo int64) (*domain.Account, error) {
	account, err := u.accountAdapter.GetAccount(ctx, accIdTo)

	if err != nil {
		u.logger.LogErrorContext(ctx, "error checking if account exist", "accIdTo", accIdTo, "error", err)
		return nil, err
	}

	if account == nil {
		u.logger.LogErrorContext(ctx, "acc to not found", "accIdTo", accIdTo)
		return nil, ErrToAccountNotFound
	}

	return account, nil
}

func (u *FundsTransferUseCaseImp) createTransaction(ctx context.Context, accFrom *domain.Account, accTo *domain.Account, amount int64) error {
	transaction := domain.NewTransaction(amount, "transfer between accounts", accFrom, accTo)

	err := u.transactionService.Save(transaction)
	if err != nil {
		u.logger.LogErrorContext(ctx, "error persisting transaction", "accFrom", accFrom, "accTo", accTo, "error", err)
		return ErrCreatingTransaction
	}

	return nil
}
