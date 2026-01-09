package core

import (
	"context"
	"fmt"

	"koyo-vpn/internal/config"
	"koyo-vpn/pkg/logger"
	"koyo-vpn/pkg/transport"

	"go.uber.org/zap"
)

type ServerCore struct {
	transport transport.GenericTransport

	logger *logger.Logger

	selfAddr string
}

func SetupServerCore(cfg *config.ServerConfig, logger *logger.Logger) (*ServerCore, error) {
	transport, err := transport.NewTransport(
		logger,
		"server",
		cfg.DeviceName,
		cfg.SelfUDPAddr,
		cfg.RemoteUDPAddr,
		cfg.Keys.Self,
		[]byte(cfg.Nonce),
		cfg.Keys.Remote,
	)
	if err != nil {
		logger.Error("failed setup transport",
			zap.Any("config", cfg),
			zap.Error(err))

		return nil, fmt.Errorf("failed setup transport: %w", err)
	}

	logger.Info("server core setupped successfully")

	return &ServerCore{
		transport: transport,
		logger:    logger,
		selfAddr:  cfg.SelfUDPAddr,
	}, nil
}

func (sc *ServerCore) Start(ctx context.Context) {
	sc.logger.Info("starting server core",
		zap.String("self_addr", sc.selfAddr))

	errors := sc.transport.StartAsync(ctx)

	for err := range errors {
		if err.Fatal {
			sc.logger.Error("fatal err from server",
				zap.String("sender", err.Sender),
				zap.String("message", err.Message))

			return
		} else {
			sc.logger.Error("err from server",
				zap.String("sender", err.Sender),
				zap.String("message", err.Message))
		}
	}
}
