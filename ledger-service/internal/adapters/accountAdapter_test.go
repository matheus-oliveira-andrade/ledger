package adapters_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/adapters"
	accountgrpc "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/adapters/grpc"
	adapters_mocks "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/adapters/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAccountAdapter_getAccount_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accId := "11"

	expectedResponse := &accountgrpc.GetAccountResponse{
		Id:        accId,
		Document:  "12345678901",
		Name:      "Test Account",
		CreatedAt: &timestamp.Timestamp{Seconds: 1234567890},
	}

	mockClient := &adapters_mocks.MockAccountClient{}
	mockClient.On("GetAccount", ctx, &accountgrpc.GetAccountRequest{AccId: accId}).Return(expectedResponse, nil)

	adapter := adapters.NewAccountAdapter(mockClient)

	accIdInt, _ := strconv.ParseInt(accId, 10, 64)

	// Act
	result, err := adapter.GetAccount(ctx, accIdInt)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	expectedAccId, _ := strconv.ParseInt(accId, 10, 64)
	assert.Equal(t, expectedAccId, result.Id)

	assert.Equal(t, "12345678901", result.Document)
	assert.Equal(t, "Test Account", result.Name)
	assert.Equal(t, expectedResponse.CreatedAt, result.CreatedAt)
	mockClient.AssertExpectations(t)
}

func TestAccountAdapter_getAccount_ClientError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accId := "11"
	expectedError := errors.New("grpc client error")

	mockClient := &adapters_mocks.MockAccountClient{}
	mockClient.On("GetAccount", ctx, &accountgrpc.GetAccountRequest{AccId: accId}).Return((*accountgrpc.GetAccountResponse)(nil), expectedError)

	adapter := adapters.NewAccountAdapter(mockClient)

	accIdInt, _ := strconv.ParseInt(accId, 10, 64)

	// Act
	result, err := adapter.GetAccount(ctx, accIdInt)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestAccountAdapter_getAccount_NilResponse(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accId := "11"

	mockClient := &adapters_mocks.MockAccountClient{}
	mockClient.On("GetAccount", ctx, &accountgrpc.GetAccountRequest{AccId: accId}).Return((*accountgrpc.GetAccountResponse)(nil), nil)

	adapter := adapters.NewAccountAdapter(mockClient)

	accIdInt, _ := strconv.ParseInt(accId, 10, 64)

	// Act
	result, err := adapter.GetAccount(ctx, accIdInt)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestAccountAdapter_getAccount_InvalidAccountId(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accId := "11"

	response := &accountgrpc.GetAccountResponse{
		Id:        "different-account-id",
		Document:  "12345678901",
		Name:      "Test Account",
		CreatedAt: &timestamp.Timestamp{Seconds: 1234567890},
	}

	mockClient := &adapters_mocks.MockAccountClient{}
	mockClient.On("GetAccount", ctx, &accountgrpc.GetAccountRequest{AccId: accId}).Return(response, nil)

	adapter := adapters.NewAccountAdapter(mockClient)

	accIdInt, _ := strconv.ParseInt(accId, 10, 64)

	// Act
	result, err := adapter.GetAccount(ctx, accIdInt)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestNewAccountAdapter(t *testing.T) {
	// Arrange
	mockClient := &adapters_mocks.MockAccountClient{}

	// Act
	adapter := adapters.NewAccountAdapter(mockClient)

	// Assert
	assert.NotNil(t, adapter)
}
