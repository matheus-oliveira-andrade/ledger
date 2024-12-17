package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/middlewares"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUseCorrelationIdMiddleware_NoCorrelationIdInRequest(t *testing.T) {

	handlerValidator := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationId := r.Context().Value(utils.CorrelationIdHeader).(string)
		assert.NotEmpty(t, correlationId)

		requestHeaderCorrelationId := r.Header.Get(string(utils.CorrelationIdHeader))
		assert.NotEmpty(t, requestHeaderCorrelationId)
	})

	middleware := middlewares.UseCorrelationIdMiddleware()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	middleware(handlerValidator).ServeHTTP(recorder, request)

	responseHeaderCorrelationid := recorder.Header().Get(string(utils.CorrelationIdHeader))
	assert.NotEmpty(t, responseHeaderCorrelationid)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestUseCorrelationIdMiddleware_CorrelationIdInRequest(t *testing.T) {
	t.Parallel()

	handlerValidator := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxCorrelationId := r.Context().Value(utils.CorrelationIdHeader).(string)
		assert.NotEmpty(t, ctxCorrelationId)

		requestHeaderCorrelationId := r.Header.Get(string(utils.CorrelationIdHeader))
		assert.NotEmpty(t, requestHeaderCorrelationId)
	})

	middleware := middlewares.UseCorrelationIdMiddleware()

	request := httptest.NewRequest(http.MethodGet, "/", nil)

	correlationId := uuid.New().String()
	request.Header.Set(string(utils.CorrelationIdHeader), correlationId)

	recorder := httptest.NewRecorder()

	middleware(handlerValidator).ServeHTTP(recorder, request)

	responseHeaderCorrelationid := recorder.Header().Get(string(utils.CorrelationIdHeader))
	assert.Equal(t, correlationId, responseHeaderCorrelationid)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
