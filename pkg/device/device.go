package device

import (
	"context"
	"fmt"
	"koyo-vpn/pkg/errors"

	"github.com/songgao/water"
)

const DefaultPacketSize = 2000

type Device struct {
	tun     *water.Interface
	fromTun chan []byte
	toTun   chan []byte
}

func NewDevice(fromTun, toTun chan []byte, name string) (*Device, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}

	config.Name = name
	tun, err := water.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed create device: %w", err)
	}

	return &Device{
		tun:     tun,
		fromTun: fromTun,
		toTun:   toTun,
	}, nil
}

func (d *Device) read() ([]byte, error) {
	buffer := make([]byte, DefaultPacketSize)
	n, err := d.tun.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

func (d *Device) write(data []byte) error {
	_, err := d.tun.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) startAsyncReader(ctx context.Context, errs chan errors.Error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			value, err := d.read()
			if err != nil {
				errs <- errors.NewError("device", err.Error(), false)
			}

			d.fromTun <- value
		}
	}
}

func (d *Device) startAsyncWriter(ctx context.Context, errs chan errors.Error) {
	for {
		select {
		case <-ctx.Done():
			close(d.fromTun)
			return
		case packet := <-d.toTun:
			if err := d.write(packet); err != nil {
				errs <- errors.NewError("device", err.Error(), false)
			}
		}
	}
}

func (d *Device) StartAsync(ctx context.Context, errors chan errors.Error) {
	go d.startAsyncReader(ctx, errors)
	go d.startAsyncWriter(ctx, errors)
}
