package middlewares_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/middlewares"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/logger"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUseLogRequestsMiddleware(t *testing.T) {
	buffer := bytes.Buffer{}
	fakeLogger := logger.NewLogger("test-service", slog.LevelInfo, &buffer, "")

	fakeHttpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := middlewares.UseLogRequestsMiddleware()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(request.Context(), utils.CtxLoggerKey, fakeLogger)

	recorder := httptest.NewRecorder()

	middleware(fakeHttpHandler).ServeHTTP(recorder, request.WithContext(ctx))

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, buffer.String(), `Request received`)
}
