package middlewares_test

import (
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/ports/http/middlewares"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/ports/http/middlewares/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseLogRequestsMiddleware(t *testing.T) {
	fakeHttpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := &middlewares_mocks.MockLogger{}
	logger.On("LogInformation", "Request received", mock.Anything).Return()

	middleware := middlewares.UseLogRequestsMiddleware(logger)

	request := httptest.NewRequest(http.MethodGet, "/", nil)

	recorder := httptest.NewRecorder()

	middleware(fakeHttpHandler).ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	logger.AssertExpectations(t)
}
