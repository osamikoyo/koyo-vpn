package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"koyo-vpn/internal/config"
	"koyo-vpn/internal/core"
	"koyo-vpn/pkg/logger"

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

	logger.Info("config loaded",
		zap.Any("config", cfg))

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
	cfg, err := config.NewConfig[*config.ServerConfig]("server", path)
	if err != nil {
		logger.Error("failed get config",
			zap.String("path", path),
			zap.Error(err))

		return nil, err
	}

	if cfg == nil {
		logger.Error("empty config",
			zap.String("path", path))

		return nil, fmt.Errorf("empty config")
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
