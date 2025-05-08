package slogger_test

import (
	"bytes"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils/slogger"
	"log/slog"
	"strings"
	"testing"
)

func TestLogger_LogInformation(t *testing.T) {
	// arrange
	var buf bytes.Buffer
	log := slogger.NewLogger("test_service", slog.LevelInfo, &buf, "")

	// assert
	log.LogInformation("message test", slog.String("name", "test"))

	// act
	output := buf.String()
	if !strings.Contains(output, `"msg":"message test"`) || !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("expected log message in output, got %v", output)
	}
}

func TestLogger_LogWarning(t *testing.T) {
	// arrange
	var buf bytes.Buffer
	log := slogger.NewLogger("test_service", slog.LevelInfo, &buf, "")

	// assert
	log.LogWarning("message test", slog.String("name", "test"))

	// act
	output := buf.String()
	if !strings.Contains(output, `"msg":"message test"`) || !strings.Contains(output, `"level":"WARN"`) {
		t.Errorf("expected log message in output, got %v", output)
	}
}

func TestLogger_LogError(t *testing.T) {
	// arrange
	var buf bytes.Buffer
	log := slogger.NewLogger("test_service", slog.LevelInfo, &buf, "")

	// assert
	log.LogError("message test", slog.String("name", "test"))

	// act
	output := buf.String()
	if !strings.Contains(output, `"msg":"message test"`) || !strings.Contains(output, `"level":"ERROR"`) {
		t.Errorf("expected log message in output, got %v", output)
	}
}

func TestLogger_MinLevel(t *testing.T) {
	// arrange
	var buf bytes.Buffer
	log := slogger.NewLogger("test_service", slog.LevelWarn, &buf, "")

	// assert
	log.LogInformation("should not appears")
	log.LogWarning("should appears 1")
	log.LogError("should appears 2")

	// act
	output := buf.String()
	if !strings.Contains(output, "should appears 1") || !strings.Contains(output, "should appears 2") {
		t.Errorf("expected log message in output, got %v", output)
	}

	if strings.Contains(output, "should not appears") {
		t.Errorf("expected log message not appear in output, got %v", output)
	}
}
