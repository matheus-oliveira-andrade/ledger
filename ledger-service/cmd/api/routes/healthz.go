package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/slogger"
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
	w.Write([]byte("OK"))

	r.logger.LogInformation("handled healthz")
}
