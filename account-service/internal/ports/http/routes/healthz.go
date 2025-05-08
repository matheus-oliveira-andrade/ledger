package routes

import (
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils/slogger"
	"net/http"

	"github.com/go-chi/chi"
)

type HealthzRoute struct {
	logger slogger.LoggerInterface
}

func NewHealthzRoute(logger slogger.LoggerInterface) *HealthzRoute {
	return &HealthzRoute{
		logger: logger,
	}
}

func (r *HealthzRoute) SetupHealthzRoutes(router chi.Router) {
	router.Get("/healthz", r.handle)
}

func (r *HealthzRoute) handle(w http.ResponseWriter, _ *http.Request) {
	r.logger.LogInformation("handling healthz")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))

	r.logger.LogInformation("handled healthz")
}
