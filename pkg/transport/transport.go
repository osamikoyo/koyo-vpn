package transport

import (
	"context"
	"fmt"
	"koyo-vpn/internal/config"
	"koyo-vpn/pkg/conn"
	"koyo-vpn/pkg/device"
	"koyo-vpn/pkg/handshake"
)

type Transport struct{
	device *device.Device
	conn *conn.Conn
	handshake *handshake.HandShakeRouter

	connInbound chan conn.Packet
	connOutbound chan conn.Packet
}

func NewTransport(deviceName, selfAddr, remoteAddr, serverKey, clientKey string) (*Transport, error) {
	var (
		fromTUN chan []byte
		toTUN chan []byte

		fromConn chan conn.Packet
		toConn chan conn.Packet
	)

	device, err := device.NewDevice(fromTUN, toTUN, deviceName)
	if err != nil{
		return nil, fmt.Errorf("failed setup device: %w", err)
	}

	conn, err := conn.NewConn(selfAddr, remoteAddr, toConn, fromConn)
	if err != nil{
		return nil, fmt.Errorf("failed setup conn: %w", err)
	}

	handshake := handshake.NewHandShakeRouter(serverKey, clientKey)

	return &Transport{
		device: device,
		conn: conn,
		handshake: handshake,
	}, nil
}

func (t *Transport) StartAsync(ctx context.Context, Ctype string) {
	
}

func (t *Transport) startServerSide(ctx context.Context, errors chan error) {
	for {
		select {
		case <- ctx.Done():
			return
		case packet := <- t.connOutbound:
			
		}
	}
}