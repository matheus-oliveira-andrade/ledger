package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/middlewares"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/routes"
	controllersV1 "github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/routes/v1"
	"github.com/matheus-oliveira-andrade/ledger/account-service/configs/settings"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/slogger"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/usecases"
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
			accountRepository := repositories.NewAccountRepository(dbConnection)

			createAccountUseCase := usecases.NewCreateAccountUseCase(logger, accountRepository)
			getAccountUseCase := usecases.NewGetAccountUseCase(logger, accountRepository)

			controllersV1.
				NewAccountsController(logger, createAccountUseCase, getAccountUseCase).
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
