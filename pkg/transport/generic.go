package transport

import "context"

type GenericTransport interface {
	StartAsync(ctx context.Context) chan error
}
