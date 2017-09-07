package redbutton

//go:generate stringer -type=Button

import (
	"errors"
	"time"

	"github.com/karalabe/hid"
)

const (
	vendor       = 0x1d34
	product      = 0x000d
	PollInterval = 200 * time.Millisecond
)

type Button int

const (
	Unknown Button = iota
	Closed
	Pressed
	Armed
)

func State(dev *hid.Device) (Button, error) {
	buf := make([]byte, 8)
	buf[0] = 0x01
	buf[7] = 0x02

	if _, err := dev.Write(buf); err != nil {
		return Unknown, err
	}

	if _, err := dev.Read(buf); err != nil {
		return Unknown, err
	}

	if buf[7] != 0x03 {
		return Unknown, nil
	}

	return Button(buf[0] & 0x03), nil
}

func Poll(dev *hid.Device, d time.Duration) <-chan Button {
	if d == 0 {
		d = PollInterval
	}
	ch := make(chan Button)
	go func() {
		prev := Unknown
		tick := time.NewTicker(d)
		defer tick.Stop()
		defer close(ch)
		for range tick.C {
			state, err := State(dev)
			if err != nil {
				return
			}
			if state != prev {
				ch <- state
				prev = state
			}
		}
	}()
	return ch
}

func Open() (*hid.Device, error) {
	devs := hid.Enumerate(vendor, product)
	if len(devs) == 0 {
		return nil, errors.New("not found")
	}
	return devs[0].Open()
}
