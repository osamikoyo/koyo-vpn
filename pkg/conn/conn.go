package conn

import (
	"context"
	"koyo-vpn/pkg/errors"
	"net"
)

type Conn struct {
	clientReady bool

	remoteSocket *net.UDPConn
	selfSocket   *net.UDPConn
	toConn       chan Packet
	fromConn     chan Packet
}

func NewConn(selfAddr, remoteAddr string, toConn, fromConn chan Packet) (*Conn, error) {
	selfUdpaddr, err := net.ResolveUDPAddr("udp", selfAddr)
	if err != nil {
		return nil, err
	}

	selfSocket, err := net.ListenUDP("udp", selfUdpaddr)
	if err != nil {
		return nil, err
	}

	return &Conn{
		toConn:     toConn,
		fromConn:   fromConn,
		selfSocket: selfSocket,
	}, nil
}

func (c *Conn) write(packet *Packet) error {
	_, err := c.remoteSocket.WriteToUDP(packet.Buf, &packet.Addr)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) read() (Packet, error) {
	buf := make([]byte, 65535)
	n, addr, err := c.selfSocket.ReadFromUDP(buf)
	if err != nil {
		return Packet{}, err
	}

	data := make([]byte, n)
	copy(data, buf[:n])

	return Packet{
		Buf:  data,
		Addr: *addr,
	}, nil
}

func (c *Conn) startAsyncReader(ctx context.Context, errs chan errors.Error, ready chan struct{}, role string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if role == "server" {
				if c.clientReady == false {
					ready <- struct{}{}

					c.clientReady = true
				}
			}

			packet, err := c.read()
			if err != nil {
				errs <- errors.NewError("conn", err.Error(), false)
			}

			c.fromConn <- packet
		}
	}
}

func (c *Conn) startAsyncWriter(ctx context.Context, errs chan errors.Error) {
	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-c.toConn:
			if err := c.write(&packet); err != nil {
				errs <- errors.NewError("conn", err.Error(), false)
			}
		}
	}
}

func (c *Conn) StartAsyncReader(ctx context.Context, errors chan errors.Error, ready chan struct{}, role string) {
	c.startAsyncReader(ctx, errors, ready, role)
}

// remote socket
func (c *Conn) StartAsyncWriter(ctx context.Context, errors chan errors.Error) {
	c.startAsyncWriter(ctx, errors)
}

func (c *Conn) SetRemoteSocket(rs *net.UDPConn) {
	c.remoteSocket = rs
}
