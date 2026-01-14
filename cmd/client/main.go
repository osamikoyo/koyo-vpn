package main

import (
	"context"
	"koyo-vpn/internal/config"
	"koyo-vpn/internal/core"
	"koyo-vpn/pkg/logger"
	"os"
	"os/signal"

	"go.uber.org/zap"
)

func main() {
	logger.Init(logger.Config{
		AppName:   "vpn-client",
		AddCaller: false,
		LogFile:   "vpn-client.log",
		LogLevel:  "debug",
	})

	logger := logger.Get()

	logger.Info("setup client...")

	cfg, err := getCfg(logger)
	if err != nil {
		return
	}

	core, err := core.SetupClientCore(cfg, logger)
	if err != nil {
		logger.Error("failed setup client core",
			zap.Error(err))
	}

	logger.Info("starting client")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	core.Start(ctx)
}

func getCfg(logger *logger.Logger) (*config.ClientConfig, error) {
	path := "client.yaml"

	for i, arg := range os.Args {
		if arg == "--config" {
			path = os.Args[i+1]
		}
	}

	cfg, err := config.NewConfig[*config.ClientConfig]("client", path)
	if err != nil {
		logger.Error("failed get config",
			zap.String("path", path),
			zap.Error(err))

		return nil, err
	}

	return cfg, nil
}
