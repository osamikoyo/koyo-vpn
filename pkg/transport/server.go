package transport

import (
	"context"
	"fmt"
	"koyo-vpn/pkg/chifer"
	"koyo-vpn/pkg/conn"
	"koyo-vpn/pkg/device"
	"koyo-vpn/pkg/errors"
	"koyo-vpn/pkg/handshake"
	"koyo-vpn/pkg/logger"
	"net"

	"go.uber.org/zap"
)

type ServerSideTransport struct {
	connInbound    chan conn.Packet
	connOutbound   chan conn.Packet
	deviceInbound  chan []byte
	deviceOutbound chan []byte

	nonce     []byte
	handshake *handshake.HandShakeRouter
	conn      *conn.Conn
	device    *device.Device
	selfAddr  *net.UDPAddr

	logger *logger.Logger
}

func newServerSideTransport(logger *logger.Logger, deviceName, selfAddr, remoteAddr, remoteKey, selfKey string, nonce []byte) (*ServerSideTransport, error) {
	var (
		fromTUN chan []byte
		toTUN   chan []byte

		fromConn chan conn.Packet
		toConn   chan conn.Packet
	)

	selfUdpAddr, err := net.ResolveUDPAddr("udp", selfAddr)
	if err != nil {
		logger.Error("failed resolve udp addr",
			zap.String("addr", selfAddr),
			zap.Error(err))

		return nil, fmt.Errorf("failed resolve udp addr %s: %w", selfAddr, err)
	}

	device, err := device.NewDevice(fromTUN, toTUN, deviceName)
	if err != nil {
		logger.Error("failed create device",
			zap.String("name", deviceName),
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

	handshake := handshake.NewHandShakeRouter(remoteKey, selfKey)

	logger.Info("transport setupped successfully",
		zap.String("device_name", deviceName),
		zap.String("self_addr", selfAddr),
		zap.String("remote_addr", remoteAddr),
		zap.String("selfKey", selfKey))

	return &ServerSideTransport{
		device:         device,
		conn:           conn,
		selfAddr:       selfUdpAddr,
		connInbound:    toConn,
		connOutbound:   fromConn,
		deviceInbound:  toTUN,
		deviceOutbound: fromTUN,
		handshake:      handshake,
	}, nil
}

func (t *ServerSideTransport) StartAsync(ctx context.Context) chan errors.Error {
	errors := make(chan errors.Error)

	go t.device.StartAsync(ctx, errors)
	go t.conn.StartAsync(ctx, errors)

	go t.startReadingFromConn(ctx, errors)
	go t.startReadingFromDevice(ctx, errors)

	t.logger.Info("transport started successfully")

	return errors
}

func (t *ServerSideTransport) startReadingFromDevice(ctx context.Context, errs chan errors.Error) {
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

func (t *ServerSideTransport) startReadingFromConn(ctx context.Context, errs chan errors.Error) {
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

func (t *ServerSideTransport) routeFromDevice(_ context.Context, buf []byte) error {
	if t.handshake.KeyIsEmpty() {
		return fmt.Errorf("empty chifer key")
	}

	key := t.handshake.GetChiferKey()

	encryptedbuf, err := chifer.Encrypt(buf, key, t.nonce)
	if err != nil {
		t.logger.Error("failed encrypt buffer",
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

func (t *ServerSideTransport) routeFromConn(_ context.Context, packet *conn.Packet) error {
	switch packet.Type {
	case 1:
		key := packet.Buf[:32]

		if err := t.handshake.NewHS(string(key)); err != nil {
			return fmt.Errorf("failed create hs: %w", err)
		}

		hspacket := conn.Packet{
			Type: 2,
			Buf:  t.handshake.GetChiferKey(),
			Addr: *t.selfAddr,
		}

		t.connInbound <- hspacket

		return nil
	case 3:
		if t.handshake.KeyIsEmpty() {
			return fmt.Errorf("empty chifer key")
		}

		key := t.handshake.GetChiferKey()

		decryptedBuf, err := chifer.Decrypt(packet.Buf, key, t.nonce)
		if err != nil {
			return err
		}

		t.deviceInbound <- decryptedBuf

		return nil
	default:
		return fmt.Errorf("unsupported packet type: %d", packet.Type)
	}
}
