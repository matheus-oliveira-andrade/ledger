package routes_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/routes"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/logger"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestHealthzHandle(t *testing.T) {
	//arrange
	buffer := bytes.Buffer{}
	fakeLogger := logger.NewLogger("test-server", slog.LevelInfo, &buffer, "")

	ctx := context.WithValue(context.Background(), utils.CtxLoggerKey, fakeLogger)

	router := chi.NewRouter()

	// act
	routes.SetupHealthz(router)

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	request = request.WithContext(ctx)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, "OK", recorder.Body.String())

	assert.Contains(t, buffer.String(), "handling healthz")
	assert.Contains(t, buffer.String(), "handled healthz")
}
