package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/slogger"
)

func UseLogRequestsMiddleware(logger slogger.LoggerInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			start := time.Now()
			next.ServeHTTP(ww, r)
			duration := time.Since(start)

			logger.LogInformation(
				"Request received",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.status),
				slog.String("remote_addr", r.RemoteAddr),
				slog.Duration("duration", duration),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
