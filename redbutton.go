package redbutton

//go:generate stringer -type=Button

import (
	"errors"
	"time"

	"github.com/karalabe/hid"
)

const (
	vendor  = 0x1d34
	product = 0x000d
)

type Button int

const (
	Unknown Button = iota
	Closed
	Pressed
	Armed
)

func State(dev *hid.Device) (Button, bool) {
	buf := make([]byte, 8)
	buf[0] = 0x01 // 0x08 ?
	buf[7] = 0x02

	if _, err := dev.Write(buf); err != nil {
		return Unknown, false
	}

	if _, err := dev.Read(buf); err != nil {
		return Unknown, false
	}

	if buf[7] != 0x03 {
		return Unknown, false
	}

	return Button(buf[0] & 0x03), true
}

func Poll(dev *hid.Device) <-chan Button {
	ch := make(chan Button)
	go func() {
		prev := Unknown
		tick := time.NewTicker(100 * time.Millisecond)
		defer tick.Stop()
		for range tick.C {
			if state, ok := State(dev); ok {
				if state != prev {
					ch <- state
				}
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
