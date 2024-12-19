package usecases_test

import (
	"errors"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/domain"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/usecases"
	usecases_mocks "github.com/matheus-oliveira-andrade/ledger/account-service/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestHandle_Err(t *testing.T) {
	// Arrange
	accId := "666"

	accountRepositoryMock := usecases_mocks.MockAccountRepository{}
	accountRepositoryMock.
		On("GetById", accId).
		Return((*domain.Account)(nil), errors.New("generic error here"))

	loggerMock := usecases_mocks.MockLogger{}
	loggerMock.On("LogInformation", "getting account by id", mock.Anything).Return()
	loggerMock.On("LogError", "error getting account by id", mock.Anything).Return()

	us := usecases.NewGetAccountUseCase(&loggerMock, &accountRepositoryMock)

	// Act
	acc, err := us.Handle(accId)

	// Assert
	assert.Nil(t, acc)
	assert.Error(t, err)

	accountRepositoryMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)

}

func TestHandle_NotFound(t *testing.T) {
	// Arrange
	accId := "666"

	accountRepositoryMock := usecases_mocks.MockAccountRepository{}
	accountRepositoryMock.
		On("GetById", accId).
		Return((*domain.Account)(nil), nil)

	loggerMock := usecases_mocks.MockLogger{}
	loggerMock.On("LogInformation", "getting account by id", mock.Anything).Return()
	loggerMock.On("LogInformation", "searched account by id", mock.Anything).Return()

	us := usecases.NewGetAccountUseCase(&loggerMock, &accountRepositoryMock)

	// Act
	acc, err := us.Handle(accId)

	// Assert
	assert.Nil(t, acc)
	assert.NoError(t, err)

	accountRepositoryMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}

func TestHandle_Success(t *testing.T) {
	// Arrange
	accId := "666"

	acc := domain.NewAccount("name test", "01234567890")
	acc.Id = accId

	accountRepositoryMock := usecases_mocks.MockAccountRepository{}
	accountRepositoryMock.
		On("GetById", accId, mock.Anything).
		Return(acc, nil)

	loggerMock := usecases_mocks.MockLogger{}
	loggerMock.On("LogInformation", "getting account by id", mock.Anything).Return()
	loggerMock.On("LogInformation", "searched account by id", mock.Anything).Return()

	us := usecases.NewGetAccountUseCase(&loggerMock, &accountRepositoryMock)

	// Act
	acc, err := us.Handle(accId)

	// Assert
	assert.NotNil(t, acc)
	assert.NoError(t, err)

	accountRepositoryMock.AssertExpectations(t)
	loggerMock.AssertExpectations(t)
}
