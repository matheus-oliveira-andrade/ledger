package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/logger"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils"
)

// ShowAccount godoc
// @Summary      Get health check
// @Description
// @Tags         healthz
// @Accept plain
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /v1/healthz [get]
func SetupHealthz(router chi.Router) {
	router.Get("/healthz", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(utils.CtxLoggerKey).(logger.LoggerInterface)

	logger.LogInformation("handling healthz")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

	logger.LogInformation("handled healthz")
}
