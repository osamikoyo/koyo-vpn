package conn

import "net"

type Packet struct {
	Buf  []byte
	Addr net.UDPAddr
	Type uint8
	// 1 is hs to server
	// 2 is hs from server
	// 3 is trafic to server
	// 4 is trafic from server
}

func ParsePacket(buf []byte, addr net.UDPAddr) *Packet {
	return &Packet{
		Buf:  buf[0:],
		Addr: addr,
		Type: buf[0],
	}
}
