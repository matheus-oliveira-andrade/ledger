package logger

import (
	"io"
	"log/slog"
	"os"
)

type LoggerInterface interface {
	LogInformation(message string, args ...any)
	LogWarning(message string, args ...any)
	LogError(message string, args ...any)
}

type Logger struct {
	slogger *slog.Logger
}

func NewLogger(serviceName string, minLevel slog.Level, output io.Writer, correlationId string) LoggerInterface {
	if output == nil {
		output = os.Stdout
	}

	jsonHandler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level:     minLevel,
		AddSource: false,
	})

	logger := slog.New(jsonHandler).
		With(slog.String("service", serviceName)).
		With(slog.String("correlationId", correlationId))

	return &Logger{
		slogger: logger,
	}
}

func (l *Logger) LogInformation(message string, args ...any) {
	l.slogger.Info(message, args...)
}

func (l *Logger) LogWarning(message string, args ...any) {
	l.slogger.Warn(message, args...)
}

func (l *Logger) LogError(message string, args ...any) {
	l.slogger.Error(message, args...)
}
