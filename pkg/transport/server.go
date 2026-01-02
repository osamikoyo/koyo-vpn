package transport

import (
	"context"
	"fmt"
	"koyo-vpn/pkg/chifer"
	"koyo-vpn/pkg/conn"
	"koyo-vpn/pkg/device"
	"koyo-vpn/pkg/handshake"
	"net"
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
}

func newServerSideTransport(deviceName, selfAddr, remoteAddr, remoteKey, selfKey string, nonce []byte) (*ServerSideTransport, error) {
	var (
		fromTUN chan []byte
		toTUN   chan []byte

		fromConn chan conn.Packet
		toConn   chan conn.Packet
	)

	selfUdpAddr, err := net.ResolveUDPAddr("udp", selfAddr)
	if err != nil {
		return nil, fmt.Errorf("failed resolve udp addr %s: %w", selfAddr, err)
	}

	device, err := device.NewDevice(fromTUN, toTUN, deviceName)
	if err != nil {
		return nil, fmt.Errorf("failed setup device: %w", err)
	}

	conn, err := conn.NewConn(selfAddr, remoteAddr, toConn, fromConn)
	if err != nil {
		return nil, fmt.Errorf("failed setup conn: %w", err)
	}

	handshake := handshake.NewHandShakeRouter(remoteKey, selfKey)

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

func (t *ServerSideTransport) StartAsync(ctx context.Context) chan error {
	errors := make(chan error)

	go t.startReadingFromConn(ctx, errors)
	go t.startReadingFromDevice(ctx, errors)

	return errors
}

func (t *ServerSideTransport) startReadingFromDevice(ctx context.Context, errors chan error) {
	for {
		select {
		case <- ctx.Done():
			return
		case buf := <- t.deviceOutbound:
			if err := t.routeFromDevice(ctx, buf);err != nil{
				errors <- err
			}
		}
	}
}

func (t *ServerSideTransport) startReadingFromConn(ctx context.Context, errors chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-t.connOutbound:
			if err := t.routeFromConn(ctx, &packet); err != nil {
				errors <- err
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
