package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/utils"
)

func UseCorrelationIdMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			correlationId := r.Header.Get(string(utils.CorrelationIdHeader))
			if correlationId == "" {
				correlationId = uuid.New().String()
			}

			ctx := context.WithValue(r.Context(), utils.CorrelationIdHeader, correlationId)

			r.Header.Set(string(utils.CorrelationIdHeader), correlationId)
			w.Header().Set(string(utils.CorrelationIdHeader), correlationId)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
