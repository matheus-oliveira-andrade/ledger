package grpc

import (
	"fmt"
	services "github.com/matheus-oliveira-andrade/ledger/account-service/internal/ports/grpc/services"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/usecases"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils/slogger"
	"google.golang.org/grpc"
	"log"
	"net"
)

type ServerWrapper struct {
	Logger slogger.LoggerInterface
	Port   int
	Server *grpc.Server
}

func NewServerWrapper(logger slogger.LoggerInterface, port int) *ServerWrapper {
	return &ServerWrapper{
		Logger: logger,
		Port:   port,
		Server: grpc.NewServer(),
	}
}

func (sw *ServerWrapper) Setup() {
	dbConnection := repositories.NewDBConnection()
	accountRepository := repositories.NewAccountRepository(dbConnection)
	getAccountUseCase := usecases.NewGetAccountUseCase(sw.Logger, accountRepository)

	accountService := services.NewAccountService(sw.Logger, getAccountUseCase)

	services.RegisterAccountServer(sw.Server, accountService)
}

func (sw *ServerWrapper) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", sw.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		sw.Logger.LogInformation("GRPC server started", "port", sw.Port)

		if err := sw.Server.Serve(listen); err != nil {
			sw.Logger.LogError("failed to serve GRPC", "error", err)
			panic("failed to serve GRPC")
		}
	}()
}
