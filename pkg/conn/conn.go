package conn

import (
	"context"
	"net"
)

type Conn struct {
	socket   *net.UDPConn
	toConn   chan Packet
	fromConn chan Packet
}

func NewConn(addr string, toConn, fromConn chan Packet) (*Conn, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		return nil, err
	}

	return &Conn{
		toConn:   toConn,
		fromConn: fromConn,
		socket:   conn,
	}, nil
}

func (c *Conn) write(packet *Packet) error {
	_, err := c.socket.WriteToUDP(packet.Buf, &packet.Addr)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) read() (Packet, error) {
	buf := make([]byte, 65535)
	n, addr, err := c.socket.ReadFromUDP(buf)
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
