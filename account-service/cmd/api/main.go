package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/matheus-oliveira-andrade/ledger/account-service/cmd/api/routes"
	"github.com/matheus-oliveira-andrade/ledger/account-service/configs/settings"
	"github.com/matheus-oliveira-andrade/ledger/account-service/configs/structuredlogs"
	"github.com/spf13/viper"
)

func main() {
	settings.Setup()
	serviceName := viper.GetString("SERVICE_NAME")
	env := viper.GetString("ENVIRONMENT")
	port := viper.GetInt("PORT")

	structuredlogs.SetupLogger(serviceName)

	r := chi.NewRouter()
	routes.SetupHealthz(r)

	slog.Info("server started", "port", port, "environment", env)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		slog.Error("server failed to start", "error", err.Error())
	}
}
