package routes

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
)

func SetupHealthz(r chi.Router) {
	r.Get("/healthz", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	slog.Info("handling healthz")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

	slog.Info("handled healthz")
}
