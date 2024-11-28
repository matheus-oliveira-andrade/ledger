package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/logger"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils"
	"github.com/spf13/viper"
)

func UseLoggerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.NewLogger(viper.GetString("SERVICE_NAME"), slog.LevelInfo, nil)
			ctx := context.WithValue(r.Context(), utils.CtxLoggerKey, logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
