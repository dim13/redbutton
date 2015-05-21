package redbutton

import (
	"errors"
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

var state = map[Button]string{
	Unknown: "Unknown",
	Closed:  "Closed",
	Pressed: "Pressed",
	Armed:   "Armed",
}

func (b Button) String() string {
	return state[b]
}

func State(dev *hid.Device) (Button, error) {
	buf := make([]byte, 8)
	buf[0] = 0x01
	buf[7] = 0x02

	if _, err := dev.Write(buf); err != nil {
		return Unknown, err
	}

	if _, err := dev.ReadTimeout(buf, 10); err != nil {
		return Unknown, err
	}

	if buf[7] != 0x03 {
		return Unknown, errors.New("bad magic")
	}

	return Button(buf[0] & 0x03), nil
}

func Poll(dev *hid.Device) <-chan Button {
	ch := make(chan Button)
	go func() {
		prev := Unknown
		for {
			state, err := State(dev)
			if err != nil {
				panic(err)
			}
			if state != prev {
				ch <- state
			}
			prev = state
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return ch
}

func Open() (*hid.Device, error) {
	return hid.Open(vendor, product, "")
}
