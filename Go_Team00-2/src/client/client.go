package main

import (
	"context"
	"net"
	"os"

	"team00/client/operator"
	"team00/client/saver"
	"team00/client/twitcher"
	"team00/config"
	"team00/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger.InitLogger("CL", config.Logger{
		Level: "debug",
	})

	ctx := context.Background()

	clientConfig := config.Client{
		Host: "127.0.0.1",
		Port: "50051",
	}

	operatorConfig := config.Operator{
		Saver: config.Saver{
			Type: saver.DBSaver,
			FileSaver: config.FileSaver{
				MeanPath:    "mean.txt",
				AnomalyPath: "anomaly.txt",
			},
			DBSaver: config.DBSaver{
				Database: "postgres",
				User:     "admin",
				Password: "admin",
			},
		},
		Calculator: config.Calculator{
			K:         0.01,
			Threshold: 20,
		},
		Threshold: 20,
	}

	op, err := operator.NewOperator(operatorConfig)
	if err != nil {
		logger.GetLogger().Err(err).Msg("failed to craete new operator")

		os.Exit(1)
	}
	creds := insecure.NewCredentials()

	conn, err := grpc.NewClient(net.JoinHostPort(clientConfig.Host, clientConfig.Port), grpc.WithTransportCredentials(creds))
	if err != nil {
		logger.GetLogger().Err(err).Msg("failed to create new connection")

		os.Exit(1)
	}

	defer conn.Close()

	tw, err := twitcher.NewTwitcher(conn, op)
	if err != nil {
		logger.GetLogger().Err(err).Msg("failed to create new twicher")

		os.Exit(1)
	}

	if err := tw.Twitch(ctx); err != nil {
		logger.GetLogger().Err(err).Msg("failed to twitch")

		os.Exit(1)
	}
}
