package routes_test

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/ports/http/routes"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/ports/http/routes/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHealthzHandle(t *testing.T) {
	// arrange
	loggerMock := &routes_mocks.MockLogger{}
	loggerMock.On("LogInformation", "handled healthz", mock.Anything).Return()
	loggerMock.On("LogInformation", "handling healthz", mock.Anything).Return()

	router := chi.NewRouter()
	routes.NewHealthzRoute(loggerMock).SetupHealthzRoutes(router)

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	// act
	router.ServeHTTP(recorder, request)

	// assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, "OK", recorder.Body.String())

	loggerMock.AssertExpectations(t)
}
