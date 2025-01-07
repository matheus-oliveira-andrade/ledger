package accountgrpc

import (
	"context"

	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/slogger"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/usecases"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	UnimplementedAccountServer

	logger            slogger.LoggerInterface
	getAccountUsecase usecases.GetAccountUseCaseInterface
}

func NewServerGRPC(
	logger slogger.LoggerInterface,
	getAccountUseCase usecases.GetAccountUseCaseInterface) *Server {
	return &Server{
		logger:            logger,
		getAccountUsecase: getAccountUseCase,
	}
}

func (s *Server) GetAccount(_ context.Context, request *GetAccountRequest) (*GetAccountResponse, error) {
	s.logger.LogInformation("searching account", "accId", request.AccId)

	acc, err := s.getAccountUsecase.Handle(request.AccId)
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
