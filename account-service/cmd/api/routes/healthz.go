package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/logger"
)

type HealthzRoute struct {
	logger logger.LoggerInterface
}

func NewHealthzRoute(logger logger.LoggerInterface) *HealthzRoute {
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
	w.Write([]byte("OK"))

	r.logger.LogInformation("handled healthz")
}
