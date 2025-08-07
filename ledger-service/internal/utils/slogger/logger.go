package slogger

import (
	"context"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/utils"
	"io"
	"log/slog"
	"os"
)

type LoggerInterface interface {
	LogInformation(message string, args ...any)
	LogWarning(message string, args ...any)
	LogError(message string, args ...any)

	LogInformationContext(ctx context.Context, message string, args ...any)
	LogWarningContext(ctx context.Context, message string, args ...any)
	LogErrorContext(ctx context.Context, message string, args ...any)
}

type Logger struct {
	slogger *slog.Logger
}

func NewLogger(serviceName string, minLevel slog.Level, output io.Writer) LoggerInterface {
	if output == nil {
		output = os.Stdout
	}

	jsonHandler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level:     minLevel,
		AddSource: false,
	})

	logger := slog.
		New(jsonHandler).
		With(slog.String("service", serviceName))

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

func (l *Logger) LogInformationContext(ctx context.Context, message string, args ...any) {
	logger, hasCorrelationId := l.withCorrelationId(ctx)
	if hasCorrelationId {
		logger.InfoContext(ctx, message, args...)
		return
	}

	l.slogger.InfoContext(ctx, message, args...)
}

func (l *Logger) LogWarningContext(ctx context.Context, message string, args ...any) {
	logger, hasCorrelationId := l.withCorrelationId(ctx)

	if hasCorrelationId {
		logger.WarnContext(ctx, message, args...)
		return
	}

	l.slogger.WarnContext(ctx, message, args...)
}

func (l *Logger) LogErrorContext(ctx context.Context, message string, args ...any) {
	logger, hasCorrelationId := l.withCorrelationId(ctx)

	if hasCorrelationId {
		logger.ErrorContext(ctx, message, args...)
		return
	}

	l.slogger.ErrorContext(ctx, message, args...)
}

func (l *Logger) withCorrelationId(ctx context.Context) (*slog.Logger, bool) {
	if ctx == nil {
		return nil, false
	}

	if correlationId, ok := ctx.Value(utils.CorrelationIdHeader).(string); ok {
		logger := l.slogger.With(slog.String("correlationId", correlationId))
		return logger, true
	}

	return nil, false
}
