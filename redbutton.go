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
	Unknown Event = iota
	Disabled
	Pressed
	Enabled
)

func Report(dev *hid.Device) (Event, error) {
	// leading zero disables sending of report number
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0, 2}
	if _, err := dev.Write(buf); err != nil {
		return Unknown, err
	}
	if _, err := dev.Read(buf[:8]); err != nil {
		return Unknown, err
	}
	return Event(buf[0] & 3), nil
}

func Poll(dev *hid.Device, d time.Duration) <-chan Event {
	if d == 0 {
		d = PollInterval
	}
	ch := make(chan Event)
	go func() {
		prev := Disabled
		tick := time.NewTicker(d)
		defer tick.Stop()
		defer close(ch)
		for range tick.C {
			ev, err := Report(dev)
			if err != nil {
				return
			}
			if ev != prev && prev != Pressed {
				ch <- ev
			}
			prev = ev
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
