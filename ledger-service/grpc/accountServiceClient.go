package accountgrpc

import (
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAccountGRPCClient() AccountClient {
	conn, err := grpc.NewClient(viper.GetString("ACCOUNT_GRPC_ADDRESS"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("GRPC conn did not connect: %v", err)
	}

	return NewAccountClient(conn)
}
