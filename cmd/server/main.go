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
	logFile, cfgPath := getPaths()

	logger.Init(logger.Config{
		AppName:   "vpn-server",
		AddCaller: false,
		LogFile:   logFile,
		LogLevel:  "debug",
	})

	logger := logger.Get()

	logger.Info("setup server...")

	cfg, err := getCfg(logger, cfgPath)
	if err != nil {
		return
	}

	core, err := core.SetupServerCore(cfg, logger)
	if err != nil {
		logger.Error("failed setup server core",
			zap.Error(err))
	}

	logger.Info("starting server")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	core.Start(ctx)
}

func getCfg(logger *logger.Logger, path string) (*config.ServerConfig, error) {
	cfg, err := config.NewConfig[*config.ServerConfig]("client", path)
	if err != nil {
		logger.Error("failed get config",
			zap.String("path", path),
			zap.Error(err))

		return nil, err
	}

	return cfg, nil
}

func getPaths() (string, string) {
	logFile := "vpn-server.log"
	cfgpath := "config.yaml"

	for i, arg := range os.Args {
		if arg == "--log-file" {
			logFile = os.Args[i+1]
		}

		if arg == "--config" {
			cfgpath = os.Args[i+1]
		}
	}


	return logFile, cfgpath
}
