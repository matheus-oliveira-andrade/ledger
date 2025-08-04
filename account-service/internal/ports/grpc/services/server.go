package accountgrpc

import (
	"context"

	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils/slogger"

	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/usecases"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountServiceImp struct {
	UnimplementedAccountServer

	logger            slogger.LoggerInterface
	getAccountUseCase usecases.GetAccountUseCaseInterface
}

func NewAccountService(
	logger slogger.LoggerInterface,
	getAccountUseCase usecases.GetAccountUseCaseInterface) *AccountServiceImp {
	return &AccountServiceImp{
		logger:            logger,
		getAccountUseCase: getAccountUseCase,
	}
}

func (s *AccountServiceImp) GetAccount(ctx context.Context, request *GetAccountRequest) (*GetAccountResponse, error) {
	s.logger.LogInformation("searching account", "accId", request.AccId)

	acc, err := s.getAccountUseCase.Handle(ctx, request.AccId)
	if err != nil {
		s.logger.LogError("error getting account", "error", err)
		return nil, err
	}

	if acc == nil {
		s.logger.LogWarning("account not found", "accId", request.AccId)
		return nil, nil
	}

	s.logger.LogInformation("searched for account", "accId", request.AccId)

	return &GetAccountResponse{
		Id:        acc.Id,
		Document:  acc.Document,
		Name:      acc.Name,
		CreatedAt: timestamppb.New(acc.CreatedAt),
	}, nil
}
