package adapters

import (
	"context"
	"strconv"

	accountgrpc "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/adapters/grpc"
	domain "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
)

type AccountAdapterInterface interface {
	GetAccount(ctx context.Context, accId int64) (*domain.Account, error)
}

type AccountAdapterImp struct {
	accountClient accountgrpc.AccountClient
}

func NewAccountAdapter(accountClient accountgrpc.AccountClient) *AccountAdapterImp {
	return &AccountAdapterImp{
		accountClient: accountClient,
	}
}

func (adapter *AccountAdapterImp) GetAccount(ctx context.Context, accId int64) (*domain.Account, error) {
	req := accountgrpc.GetAccountRequest{
		AccId: strconv.FormatInt(accId, 10),
	}

	acc, err := adapter.accountClient.GetAccount(ctx, &req)
	if err != nil {
		return nil, err
	}

	if acc == nil || acc.Id != req.AccId {
		return nil, nil
	}

	return domain.NewAccount(accId, acc.Document, acc.Name, acc.CreatedAt), nil
}
