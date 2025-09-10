package http

import (
	"fmt"
	"net/http"

	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/adapters"
	accountgrpc "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/adapters/grpc"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/ports/http/middlewares"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/ports/http/routes"
	controllersV1 "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/ports/http/routes/v1"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/utils/slogger"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/services"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/usecases"
)

type ServerWrapper struct {
	ServiceName string
	Logger      slogger.LoggerInterface
	Port        int
	Environment string
	Router      *chi.Mux
}

func NewServerWrapper(serviceName string, logger slogger.LoggerInterface, port int, environment string) *ServerWrapper {
	return &ServerWrapper{
		ServiceName: serviceName,
		Logger:      logger,
		Port:        port,
		Environment: environment,
		Router:      chi.NewRouter(),
	}
}

func (hs *ServerWrapper) Setup() {
	hs.Router.Use(middlewares.UseCorrelationIdMiddleware())
	hs.Router.Use(middlewares.UseLogRequestsMiddleware(hs.Logger))
	hs.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	hs.Router.Route("/api", func(apiRouter chi.Router) {
		routes.NewHealthzRoute(hs.Logger).SetupHealthzRoutes(apiRouter)

		apiRouter.Route("/v1", func(v1Router chi.Router) {
			dbConnection := repositories.NewDBConnection()

			transactionRepository := repositories.NewTransactionRepository(dbConnection)
			transactionLineRepository := repositories.NewTransactionLineRepository(dbConnection)
			transactionService := services.NewTransactionService(hs.Logger, transactionRepository, transactionLineRepository)
			balanceService := services.NewBalanceService(transactionLineRepository)

			accountClient := accountgrpc.NewAccountGRPCClient()
			accountAdapter := adapters.NewAccountAdapter(accountClient)

			fundsTransferUseCase := usecases.NewFundsTransferUseCase(hs.Logger, transactionService, accountAdapter, balanceService)

			controllersV1.NewFundsTransferController(hs.Logger, fundsTransferUseCase).
				RegisterRoutes(v1Router)

			getBalanceUseCase := usecases.NewGetBalanceUseCase(hs.Logger, balanceService)

			controllersV1.NewBalanceController(hs.Logger, getBalanceUseCase).
				RegisterRoutes(v1Router)

			statementRepository := repositories.NewStatementRepository(dbConnection)
			getStatementUseCase := usecases.NewGetStatementUseCase(hs.Logger, statementRepository)

			controllersV1.NewStatementController(hs.Logger, getStatementUseCase).
				RegisterRoutes(v1Router)
		})
	})
}

func (hs *ServerWrapper) Start() {
	startServer(hs.Logger, hs.Port, hs.Environment, hs.Router)
}

func startServer(l slogger.LoggerInterface, port int, env string, r *chi.Mux) {
	l.LogInformation("server started", "port", port, "environment", env)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		l.LogError("server failed to start", "error", err.Error())
	}
}
