package redbutton

//go:generate stringer -type=Button

import (
	"time"

	"github.com/GeertJohan/go.hid"
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
	buf[0] = 0x01
	buf[7] = 0x02

	if _, err := dev.Write(buf); err != nil {
		return Unknown, false
	}

	if _, err := dev.ReadTimeout(buf, 10); err != nil {
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
		for {
			if state, ok := State(dev); ok {
				if state != prev {
					ch <- state
				}
				prev = state
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return ch
}

func Open() (*hid.Device, error) {
	return hid.Open(vendor, product, "")
}
