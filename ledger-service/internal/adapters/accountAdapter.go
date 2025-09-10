package adapters

import (
	"context"

	accountgrpc "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/adapters/grpc"
	domain "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
)

type AccountAdapterInterface interface {
	GetAccount(ctx context.Context, accId string) (*domain.Account, error)
}

type AccountAdapterImp struct {
	accountClient accountgrpc.AccountClient
}

func NewAccountAdapter(accountClient accountgrpc.AccountClient) *AccountAdapterImp {
	return &AccountAdapterImp{
		accountClient: accountClient,
	}
}

func (adapter *AccountAdapterImp) GetAccount(ctx context.Context, accId string) (*domain.Account, error) {
	req := accountgrpc.GetAccountRequest{
		AccId: accId,
	}

	acc, err := adapter.accountClient.GetAccount(ctx, &req)
	if err != nil {
		return nil, err
	}

	if acc == nil || acc.Id != req.AccId {
		return nil, nil
	}

	return domain.NewAccount(acc.Id, acc.Document, acc.Name, acc.CreatedAt), nil
}
