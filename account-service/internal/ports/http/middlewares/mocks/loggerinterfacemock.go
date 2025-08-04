package middlewares_mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) LogInformation(message string, args ...any) {
	m.Called(message, args)
}

func (m *MockLogger) LogWarning(message string, args ...any) {
	m.Called(message, args)
}

func (m *MockLogger) LogError(message string, args ...any) {
	m.Called(message, args)
}

func (m *MockLogger) LogInformationContext(ctx context.Context, message string, args ...any) {
	m.Called(ctx, message, args)
}

func (m *MockLogger) LogWarningContext(ctx context.Context, message string, args ...any) {
	m.Called(ctx, message, args)
}

func (m *MockLogger) LogErrorContext(ctx context.Context, message string, args ...any) {
	m.Called(ctx, message, args)
}
