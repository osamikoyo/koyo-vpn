package conn

import "net"

type Packet struct {
    Buf      []byte         
    Addr     net.UDPAddr    
    Type     uint8         
}

func ParsePacket(buf []byte, addr net.UDPAddr) *Packet {
    return &Packet{
        Buf: buf[0:],
        Addr: addr,
        Type: buf[0],
    }
}