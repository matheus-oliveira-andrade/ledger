package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/middlewares"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/routes"
	"github.com/matheus-oliveira-andrade/ledger/account-service/configs/settings"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/logger"
	"github.com/spf13/viper"
)

func main() {
	settings.Setup()
	serviceName := viper.GetString("SERVICE_NAME")
	if serviceName == "" {
		panic("env variable SERVICE_NAME not loaded")
	}

	env := viper.GetString("ENVIRONMENT")
	if env == "" {
		panic("env variable ENVIRONMENT not loaded")
	}

	port := viper.GetInt("PORT")
	if port == 0 {
		panic("env variable PORT not loaded")
	}

	r := chi.NewRouter()
	r.Use(middlewares.UseLoggerMiddleware())
	r.Use(middlewares.UseRequestLoggerMiddleware())

	routes.SetupHealthz(r)

	logger := logger.NewLogger(serviceName, slog.LevelInfo, nil)
	logger.LogInformation("server started", "port", port, "environment", env)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		logger.LogError("server failed to start", "error", err.Error())
	}
}
