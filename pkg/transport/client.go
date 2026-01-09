package transport

import (
	"context"
	"fmt"
	"koyo-vpn/pkg/chifer"
	"koyo-vpn/pkg/conn"
	"koyo-vpn/pkg/device"
	"koyo-vpn/pkg/errors"
	"koyo-vpn/pkg/logger"
	"net"

	"go.uber.org/zap"
)

type ClientSideTransport struct {
	connInbound    chan conn.Packet
	connOutbound   chan conn.Packet
	deviceInbound  chan []byte
	deviceOutbound chan []byte

	nonce    []byte
	conn     *conn.Conn
	device   *device.Device
	selfAddr *net.UDPAddr
	key      []byte

	logger *logger.Logger
}

func newClientSideTransport(logger *logger.Logger, deviceName, selfAddr, remoteAddr, selfKey string, nonce []byte) (*ClientSideTransport, error) {
	var (
		fromTUN chan []byte
		toTUN   chan []byte

		fromConn chan conn.Packet
		toConn   chan conn.Packet
	)

	selfUdpAddr, err := net.ResolveUDPAddr("udp", selfAddr)
	if err != nil {
		logger.Error("faield setup self udp addr",
			zap.String("addr", selfAddr),
			zap.Error(err))

		return nil, fmt.Errorf("failed resolve udp addr %s: %w", selfAddr, err)
	}

	device, err := device.NewDevice(fromTUN, toTUN, deviceName)
	if err != nil {
		logger.Error("failed setup device",
			zap.String("device_name", deviceName),
			zap.Error(err))

		return nil, fmt.Errorf("failed setup device: %w", err)
	}

	conn, err := conn.NewConn(selfAddr, remoteAddr, toConn, fromConn)
	if err != nil {
		logger.Error("failed setup conn",
			zap.String("self_addr", selfAddr),
			zap.String("remote_addr", remoteAddr),
			zap.Error(err))

		return nil, fmt.Errorf("failed setup conn: %w", err)
	}

	return &ClientSideTransport{
		device:         device,
		conn:           conn,
		selfAddr:       selfUdpAddr,
		connInbound:    toConn,
		connOutbound:   fromConn,
		deviceInbound:  toTUN,
		deviceOutbound: fromTUN,
		logger:         logger,
	}, nil
}

func (t *ClientSideTransport) StartAsync(ctx context.Context) chan errors.Error {
	errors := make(chan errors.Error, 15)

	go t.device.StartAsync(ctx, errors)
	go t.conn.StartAsync(ctx, errors)

	go t.startReadingFromConn(ctx, errors)
	go t.startReadingFromDevice(ctx, errors)

	return errors
}

func (t *ClientSideTransport) startReadingFromDevice(ctx context.Context, errs chan errors.Error) {
	for {
		select {
		case <-ctx.Done():
			return
		case buf := <-t.deviceOutbound:
			if err := t.routeFromDevice(ctx, buf); err != nil {
				errs <- errors.NewError("transport", err.Error(), false)
			}
		}
	}
}

func (t *ClientSideTransport) startReadingFromConn(ctx context.Context, errs chan errors.Error) {
	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-t.connOutbound:
			if err := t.routeFromConn(ctx, &packet); err != nil {
				errs <- errors.NewError("transport", err.Error(), false)
			}
		}
	}
}

func (t *ClientSideTransport) routeFromDevice(_ context.Context, buf []byte) error {
	if len(t.key) == 0 {
		t.logger.Error("empty chifer key")

		return fmt.Errorf("empty chifer key")
	}

	encryptedbuf, err := chifer.Encrypt(buf, t.key, t.nonce)
	if err != nil {
		t.logger.Error("failed encrypt",
			zap.ByteString("buf", buf),
			zap.ByteString("key", t.key),
			zap.Error(err))

		return fmt.Errorf("failed ecrypt buffer: %w", err)
	}

	pType := 1

	p := conn.Packet{
		Addr: *t.selfAddr,
		Buf:  encryptedbuf,
		Type: uint8(pType),
	}

	t.connInbound <- p

	return nil
}

func (t *ClientSideTransport) routeFromConn(_ context.Context, packet *conn.Packet) error {
	switch packet.Type {
	case 2:
		key := packet.Buf[:32]

		t.key = key

		return nil
	case 4:
		if len(t.key) == 0 {
			return fmt.Errorf("empty chifer key")
		}

		decryptedBuf, err := chifer.Decrypt(packet.Buf, []byte(t.key), t.nonce)
		if err != nil {
			return err
		}

		t.deviceInbound <- decryptedBuf

		return nil
	default:
		t.logger.Error("unsupported packet type",
			zap.Any("packet", packet))

		return fmt.Errorf("unsupported packet type: %d", packet.Type)
	}
}
