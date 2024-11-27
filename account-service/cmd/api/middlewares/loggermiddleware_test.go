package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/middlewares"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/utils"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUseLoggerMiddleware(t *testing.T) {
	viper.Set("SERVICE_NAME", "test-service")

	handlerValidator := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logValue := r.Context().Value(utils.CtxLoggerKey).(logger.LoggerInterface)

		assert.NotNil(t, logValue, "logger not found in context")

		w.WriteHeader(http.StatusOK)
	})

	middleware := middlewares.UseLoggerMiddleware()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	middleware(handlerValidator).ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
