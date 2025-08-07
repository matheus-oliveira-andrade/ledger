package main

import (
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/configs/settings"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/ports/http"
	"github.com/matheus-oliveira-andrade/ledger/ledger-service/internal/utils/slogger"
	"github.com/spf13/viper"
	"log/slog"
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

	logger := slogger.NewLogger(serviceName, slog.LevelInfo, nil)

	httpServer := http.NewServerWrapper(serviceName, logger, port, env)
	httpServer.Setup()
	httpServer.Start()

	logger.LogInformation("Server started", "httpPort", port, "env", env)
}
