package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/middlewares"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/routes/v1"
	"github.com/matheus-oliveira-andrade/ledger/account-service/configs/settings"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/logger"
	"github.com/spf13/viper"

	_ "github.com/matheus-oliveira-andrade/ledger/account-service/docs"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
)

// @title Account service Swagger API
// @version 1.0
// @description account service part of ledger project
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
	r.Use(middlewares.UseLogRequestsMiddleware())

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:"+strconv.Itoa(port)+"/swagger/doc.json"),
	))

	r.Route("/v1", func(r chi.Router) {
		routes.SetupHealthz(r)
	})

	logger := logger.NewLogger(serviceName, slog.LevelInfo, nil, uuid.NewString())
	logger.LogInformation("server started", "port", port, "environment", env)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		logger.LogError("server failed to start", "error", err.Error())
	}
}
