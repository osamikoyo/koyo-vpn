package transport

import (
	"koyo-vpn/pkg/conn"
	"koyo-vpn/pkg/device"
)

type Transport struct{
	Device *device.Device
	Conn *conn.Conn
}

func NewTransport()