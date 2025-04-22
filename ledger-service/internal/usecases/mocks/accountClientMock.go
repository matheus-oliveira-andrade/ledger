package usecases_mocks

import (
	"context"
	accountgrpc "github.com/matheus-oliveira-andrade/ledger/ledger-service/grpc"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAccountClient struct {
	mock.Mock
}

func (m *MockAccountClient) GetAccount(ctx context.Context, in *accountgrpc.GetAccountRequest, _ ...grpc.CallOption) (*accountgrpc.GetAccountResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*accountgrpc.GetAccountResponse), args.Error(1)
}
