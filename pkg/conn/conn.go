package conn

import (
	"context"
	"net"
)

type Conn struct {
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

	remoteUDPAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		return nil, err
	}

	remoteSocket, err := net.DialUDP("udp", nil, remoteUDPAddr)
	if err != nil {
		return nil, err
	}

	return &Conn{
		toConn:       toConn,
		fromConn:     fromConn,
		selfSocket:   selfSocket,
		remoteSocket: remoteSocket,
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

func (c *Conn) startAsyncReader(ctx context.Context, errors chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			packet, err := c.read()
			if err != nil {
				select {
				case errors <- err:
				default:
				}
			}

			c.fromConn <- packet
		}
	}
}

func (c *Conn) startAsyncWriter(ctx context.Context, errors chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-c.toConn:
			if err := c.write(&packet); err != nil {
				select {
				case errors <- err:
				default:
				}
			}
		}
	}
}

func (c *Conn) StartAsync(ctx context.Context) chan error {
	errors := make(chan error, 5)

	go c.startAsyncReader(ctx, errors)
	go c.startAsyncWriter(ctx, errors)

	return errors
}
