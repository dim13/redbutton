package redbutton

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

//go:generate stringer -type=Event

type Event int

const (
	Unknown       Event = 0x00
	LidClosed     Event = 0x15
	ButtonPressed Event = 0x16
	LidOpen       Event = 0x17
)

func State(dev *hid.Device) (Event, error) {
	// leading zero disables sending of report number
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0, 2}
	if _, err := dev.Write(buf); err != nil {
		return Unknown, err
	}
	if _, err := dev.Read(buf[:8]); err != nil {
		return Unknown, err
	}
	return Event(buf[0]), nil
}

func Poll(dev *hid.Device, d time.Duration) <-chan Event {
	if d == 0 {
		d = PollInterval
	}
	ch := make(chan Event)
	go func() {
		prev := LidClosed
		tick := time.NewTicker(d)
		defer tick.Stop()
		defer close(ch)
		for range tick.C {
			state, err := State(dev)
			if err != nil {
				return
			}
			if state != prev && prev != ButtonPressed {
				ch <- state
			}
			prev = state
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
