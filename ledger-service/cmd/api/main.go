package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/cmd/api/middlewares"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/cmd/api/routes"
	controllersV1 "github.com/matheus-oliveira-andrade/ledger/ledger-service/cmd/api/routes/v1"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/configs/settings"
	accountgrpc "github.com/matheus-oliveira-andrade/ledger/ledger-service/grpc"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/services"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/slogger"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/usecases"
	"github.com/spf13/viper"
)

func main() {
	settings.Setup()
	serviceName := viper.GetString("SERVICE_NAME")
	if serviceName == "" {
		panic("env variable SERVICE_NAME not loaded")
	}

	env := viper.GetString("ENVIRONMENT")
	if env == "" {
		panic("env variable ENVIRONMENT not loaded")
	}

	port := viper.GetInt("PORT")
	if port == 0 {
		panic("env variable PORT not loaded")
	}

	logger := slogger.NewLogger(serviceName, slog.LevelInfo, nil, uuid.NewString())

	r := chi.NewRouter()
	r.Use(middlewares.UseCorrelationIdMiddleware())
	r.Use(middlewares.UseLogRequestsMiddleware(logger))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	r.Route("/api", func(apiRouter chi.Router) {
		routes.NewHealthzRoute(logger).SetupHealthzRoutes(apiRouter)

		apiRouter.Route("/v1", func(v1Router chi.Router) {
			dbConnection := repositories.NewDBConnection()

			transactionRepository := repositories.NewTransactionRepository(dbConnection)
			transactionLineRepository := repositories.NewTransactionLineRepository(dbConnection)
			transactionService := services.NewTransactionService(logger, transactionRepository, transactionLineRepository)
			balanceService := services.NewBalanceService(transactionLineRepository)

			accountClient := accountgrpc.NewAccountGRPCClient()

			fundsTransferUseCase := usecases.NewFundsTransferUseCase(logger, transactionService, accountClient, balanceService)

			controllersV1.
				NewFundsTransferController(logger, fundsTransferUseCase).
				RegisterRoutes(v1Router)

			getBalanceUseCase := usecases.NewGetBalanceUseCase(logger, balanceService)

			controllersV1.
				NewBalanceController(logger, getBalanceUseCase).
				RegisterRoutes(v1Router)

			statementRepository := repositories.NewStatementRepository(dbConnection)
			getStatementUseCase := usecases.NewGetStatementUseCase(logger, statementRepository)

			controllersV1.
				NewStatementController(logger, getStatementUseCase).
				RegisterRoutes(v1Router)
		})
	})

	startServer(logger, port, env, r)
}

func startServer(l slogger.LoggerInterface, port int, env string, r *chi.Mux) {
	l.LogInformation("server started", "port", port, "environment", env)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		l.LogError("server failed to start", "error", err.Error())
	}
}
