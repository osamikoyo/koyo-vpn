package transport

import (
	"context"
	"fmt"

	"koyo-vpn/pkg/errors"
	"koyo-vpn/pkg/logger"
)

type GenericTransport interface {
	StartAsync(ctx context.Context) chan errors.Error
}

func NewTransport(
	logger *logger.Logger,
	side string,
	deviceName string,
	selfAddr string,
	remoteAddr string,
	selfKey string,
	remoteKey ...string,
) (GenericTransport, error) {
	switch side {
	case "server":
		if len(remoteKey) == 0 {
			return nil, fmt.Errorf("remote key is nil")
		}

		return newServerSideTransport(
			logger,
			deviceName,
			selfAddr,
			remoteAddr,
			remoteKey[0],
			selfKey,
		)
	case "client":
		return newClientSideTransport(
			logger,
			deviceName,
			selfAddr,
			remoteAddr,
			selfKey,
		)
	default:
		return nil, fmt.Errorf("unsupported side: %s", side)
	}
}
