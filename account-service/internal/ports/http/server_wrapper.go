package http

import (
	"fmt"
	middlewares "github.com/matheus-oliveira-andrade/ledger/account-service/internal/ports/http/middlewares"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/ports/http/routes"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/ports/http/routes/v1"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils/slogger"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/repositories"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/usecases"
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
			accountRepository := repositories.NewAccountRepository(dbConnection)

			createAccountUseCase := usecases.NewCreateAccountUseCase(hs.Logger, accountRepository)
			getAccountUseCase := usecases.NewGetAccountUseCase(hs.Logger, accountRepository)

			controllersV1.NewAccountsController(hs.Logger, createAccountUseCase, getAccountUseCase).
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
