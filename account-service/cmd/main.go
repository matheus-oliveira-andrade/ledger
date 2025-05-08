package main

import (
	"github.com/google/uuid"
	"github.com/matheus-oliveira-andrade/ledger/account-service/configs/settings"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/ports/grpc"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/ports/http"
	"github.com/matheus-oliveira-andrade/ledger/account-service/internal/utils/slogger"
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

	httpPort := viper.GetInt("PORT")
	if httpPort == 0 {
		panic("env variable PORT not loaded")
	}

	rpcPort := viper.GetInt("RPC_PORT")
	if rpcPort == 0 {
		panic("env variable RPC_PORT not loaded")
	}

	logger := slogger.NewLogger(serviceName, slog.LevelInfo, nil, uuid.NewString())

	grpcServer := grpc.NewServerWrapper(logger, rpcPort)
	grpcServer.Setup()
	grpcServer.Start()

	httpServer := http.NewServerWrapper(serviceName, logger, httpPort, env)
	httpServer.Setup()
	httpServer.Start()

	logger.LogInformation("Servers started", "httpPort", httpPort, "rpcPort", rpcPort, "env", env)
}
