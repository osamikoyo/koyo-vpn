package core

import (
	"context"
	"fmt"
	"koyo-vpn/internal/config"
	"koyo-vpn/pkg/logger"
	"koyo-vpn/pkg/transport"

	"go.uber.org/zap"
)

type ClientCore struct {
	transport transport.GenericTransport

	selfAddr string

	logger *logger.Logger
}

func SetupClientCore(cfg *config.ClientConfig, logger *logger.Logger) (*ClientCore, error) {
	transport, err := transport.NewTransport(logger, "client", cfg.DeviceName, cfg.Addrs.Self, cfg.Addrs.Remote, cfg.SelfKey)
	if err != nil {
		logger.Error("failed setup transport",
			zap.Any("cfg", cfg),
			zap.Error(err))

		return nil, fmt.Errorf("failed setup transport: %w", err)
	}

	logger.Info("client core setupped successfully")

	return &ClientCore{
		transport: transport,
		logger:    logger,
	}, nil
}

func (cc *ClientCore) Start(ctx context.Context) {
	cc.logger.Info("starting client core",
		zap.String("self_addr", cc.selfAddr))

	errors := cc.transport.StartAsync(ctx)

	for err := range errors {
		if err.Fatal {
			cc.logger.Error("fatal err from server",
				zap.String("sender", err.Sender),
				zap.String("message", err.Message))

			return
		} else {
			cc.logger.Error("err from server",
				zap.String("sender", err.Sender),
				zap.String("message", err.Message))
		}
	}
}
