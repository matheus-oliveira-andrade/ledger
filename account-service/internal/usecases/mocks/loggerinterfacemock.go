package usecases_mocks

import (
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
