package services_test

import (
	"errors"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/services"
	services_mocks "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBalanceService_CalculateBalance_Success(t *testing.T) {
	// Arrange
	accId := int64(12323)

	transactionLineMock := services_mocks.MockTransactionLineRepository{}
	transactionLineMock.
		On("GetTransactions", accId).
		Return(&[]domain.TransactionLine{
			*domain.NewTransactionLine(accId, 100, domain.Credit),
			*domain.NewTransactionLine(accId, 100, domain.Credit),
		}, nil)

	service := services.NewBalanceService(&transactionLineMock)

	// Act
	result, err := service.CalculateBalance(accId)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(200), result)
}

func TestBalanceService_CalculateBalance_Error(t *testing.T) {
	// Arrange
	accId := int64(12323)

	expectedError := errors.New("random error")

	transactionLineMock := services_mocks.MockTransactionLineRepository{}
	transactionLineMock.
		On("GetTransactions", accId).
		Return(&[]domain.TransactionLine{}, expectedError)

	service := services.NewBalanceService(&transactionLineMock)

	// Act
	result, err := service.CalculateBalance(accId)

	// Assert
	assert.Error(t, expectedError, err)
	assert.Equal(t, int64(0), result)
}
