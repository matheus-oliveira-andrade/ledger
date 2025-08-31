package usecases_test

import (
	"context"
	"errors"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/adapters/grpc"
	usecases_mocks "github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/usecases/mocks"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/utils/slogger"
	"strconv"
	"testing"

	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testMocks struct {
	mockLogger             *slogger.MockLogger
	mockTransactionService *usecases_mocks.MockTransactionService
	mockAccountClient      *usecases_mocks.MockAccountClient
	mockBalanceService     *usecases_mocks.MockBalanceService
}

func (tm *testMocks) setup() {
	tm.mockLogger = &slogger.MockLogger{}
	tm.mockLogger.
		On("LogInformationContext", mock.Anything, mock.Anything, mock.Anything).
		Return()
	tm.mockLogger.
		On("LogErrorContext", mock.Anything, mock.Anything, mock.Anything).
		Return()

	tm.mockTransactionService = &usecases_mocks.MockTransactionService{}
	tm.mockAccountClient = &usecases_mocks.MockAccountClient{}
	tm.mockBalanceService = &usecases_mocks.MockBalanceService{}
}

func newTestMocks() *testMocks {
	tm := &testMocks{}
	tm.setup()

	return tm
}

func TestFundsTransferUseCase_Handle_ErrorCheckingFromAccountExist(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accFrom := int64(1)
	accTo := int64(2)
	amount := int64(100)

	testMocks := newTestMocks()

	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.Anything).
		Return((*accountgrpc.GetAccountResponse)(nil), errors.New("connection error"))

	useCase := usecases.NewFundsTransferUseCase(
		testMocks.mockLogger,
		testMocks.mockTransactionService,
		testMocks.mockAccountClient,
		testMocks.mockBalanceService,
	)

	// Act
	err := useCase.Handle(ctx, accFrom, accTo, amount)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "connection error", err.Error())
	testMocks.mockAccountClient.AssertExpectations(t)
}

func TestFundsTransferUseCase_Handle_FromAccountNotExist(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accFrom := int64(1)
	accTo := int64(2)
	amount := int64(100)

	testMocks := newTestMocks()

	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.Anything).
		Return((*accountgrpc.GetAccountResponse)(nil), nil)

	useCase := usecases.NewFundsTransferUseCase(
		testMocks.mockLogger,
		testMocks.mockTransactionService,
		testMocks.mockAccountClient,
		testMocks.mockBalanceService,
	)

	// Act
	err := useCase.Handle(ctx, accFrom, accTo, amount)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "acc from not found", err.Error())
	testMocks.mockAccountClient.AssertExpectations(t)
}

func TestFundsTransferUseCase_Handle_ErrorGettingFromAccountBalance(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accFrom := int64(1)
	accTo := int64(2)
	amount := int64(100)

	testMocks := newTestMocks()

	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.Anything).
		Return(&accountgrpc.GetAccountResponse{Id: "1"}, nil)

	testMocks.mockBalanceService.
		On("CalculateBalance", ctx, accFrom).
		Return(int64(0), errors.New("balance error"))

	useCase := usecases.NewFundsTransferUseCase(
		testMocks.mockLogger,
		testMocks.mockTransactionService,
		testMocks.mockAccountClient,
		testMocks.mockBalanceService,
	)

	// Act
	err := useCase.Handle(ctx, accFrom, accTo, amount)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "balance error", err.Error())
	testMocks.mockAccountClient.AssertExpectations(t)
	testMocks.mockBalanceService.AssertExpectations(t)
}

func TestFundsTransferUseCase_Handle_FromAccountInsufficientBalance(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accFrom := int64(1)
	accTo := int64(2)
	amount := int64(100)

	testMocks := newTestMocks()

	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.Anything).
		Return(&accountgrpc.GetAccountResponse{Id: "1"}, nil)

	testMocks.mockBalanceService.
		On("CalculateBalance", ctx, accFrom).
		Return(int64(50), nil)

	useCase := usecases.NewFundsTransferUseCase(
		testMocks.mockLogger,
		testMocks.mockTransactionService,
		testMocks.mockAccountClient,
		testMocks.mockBalanceService,
	)

	// Act
	err := useCase.Handle(ctx, accFrom, accTo, amount)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "insufficient balance", err.Error())
	testMocks.mockAccountClient.AssertExpectations(t)
	testMocks.mockBalanceService.AssertExpectations(t)
}

func TestFundsTransferUseCase_Handle_ErrorCheckingToAccountExist(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accFrom := int64(1)
	accTo := int64(2)
	amount := int64(100)

	testMocks := newTestMocks()

	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.Anything).
		Return(&accountgrpc.GetAccountResponse{Id: "1"}, nil).Once()
	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.Anything).
		Return((*accountgrpc.GetAccountResponse)(nil), errors.New("connection error")).Once()

	testMocks.mockBalanceService.
		On("CalculateBalance", ctx, accFrom).
		Return(int64(200), nil)

	useCase := usecases.NewFundsTransferUseCase(
		testMocks.mockLogger,
		testMocks.mockTransactionService,
		testMocks.mockAccountClient,
		testMocks.mockBalanceService,
	)

	// Act
	err := useCase.Handle(ctx, accFrom, accTo, amount)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "connection error", err.Error())
	testMocks.mockAccountClient.AssertExpectations(t)
	testMocks.mockBalanceService.AssertExpectations(t)
}

func TestFundsTransferUseCase_Handle_ToAccountNotExist(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accFrom := int64(1)
	accTo := int64(2)
	amount := int64(100)

	testMocks := newTestMocks()

	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.Anything).
		Return(&accountgrpc.GetAccountResponse{Id: "1"}, nil).Once()
	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.Anything).
		Return((*accountgrpc.GetAccountResponse)(nil), nil).Once()

	testMocks.mockBalanceService.
		On("CalculateBalance", ctx, accFrom).
		Return(int64(200), nil)

	useCase := usecases.NewFundsTransferUseCase(
		testMocks.mockLogger,
		testMocks.mockTransactionService,
		testMocks.mockAccountClient,
		testMocks.mockBalanceService,
	)

	// Act
	err := useCase.Handle(ctx, accFrom, accTo, amount)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "acc to not found", err.Error())
	testMocks.mockAccountClient.AssertExpectations(t)
	testMocks.mockBalanceService.AssertExpectations(t)
}

func TestFundsTransferUseCase_Handle_ErrorSaveTransaction(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accFrom := int64(1)
	accTo := int64(2)
	amount := int64(100)

	testMocks := newTestMocks()

	testMocks.mockTransactionService.
		On("Save", mock.Anything).
		Return(errors.New("save error"))

	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.Anything).
		Return(&accountgrpc.GetAccountResponse{Id: "1"}, nil).Once()
	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.Anything).
		Return(&accountgrpc.GetAccountResponse{Id: "2"}, nil).Once()

	testMocks.mockBalanceService.
		On("CalculateBalance", ctx, accFrom).
		Return(int64(200), nil)

	useCase := usecases.NewFundsTransferUseCase(
		testMocks.mockLogger,
		testMocks.mockTransactionService,
		testMocks.mockAccountClient,
		testMocks.mockBalanceService,
	)

	// Act
	err := useCase.Handle(ctx, accFrom, accTo, amount)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "error persisting transaction", err.Error())
	testMocks.mockAccountClient.AssertExpectations(t)
	testMocks.mockBalanceService.AssertExpectations(t)
	testMocks.mockTransactionService.AssertExpectations(t)
}

func TestFundsTransferUseCase_Handle_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	accFrom := int64(1)
	accTo := int64(2)
	amount := int64(100)

	testMocks := newTestMocks()

	testMocks.mockTransactionService.
		On("Save", mock.Anything).
		Return(nil)

	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.MatchedBy(func(req *accountgrpc.GetAccountRequest) bool {
			return req.AccId == strconv.FormatInt(accFrom, 10)
		}), mock.Anything).
		Return(&accountgrpc.GetAccountResponse{Id: "1"}, nil).Once()
	testMocks.mockAccountClient.
		On("GetAccount", ctx, mock.MatchedBy(func(req *accountgrpc.GetAccountRequest) bool {
			return req.AccId == strconv.FormatInt(accTo, 10)
		}), mock.Anything).
		Return(&accountgrpc.GetAccountResponse{Id: "2"}, nil).Once()

	testMocks.mockBalanceService.
		On("CalculateBalance", ctx, accFrom).
		Return(int64(200), nil)

	useCase := usecases.NewFundsTransferUseCase(
		testMocks.mockLogger,
		testMocks.mockTransactionService,
		testMocks.mockAccountClient,
		testMocks.mockBalanceService,
	)

	// Act
	err := useCase.Handle(ctx, accFrom, accTo, amount)

	// Assert
	assert.NoError(t, err)
	testMocks.mockAccountClient.AssertExpectations(t)
	testMocks.mockBalanceService.AssertExpectations(t)
	testMocks.mockTransactionService.AssertExpectations(t)
	testMocks.mockLogger.AssertCalled(t, "LogInformationContext", mock.Anything, "transfer completed", mock.Anything)
}
