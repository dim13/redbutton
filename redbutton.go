package redbutton

import (
	"errors"
	"time"

	"github.com/GeertJohan/go.hid"
)

type Button struct {
	Button bool
	Lid    bool
}

func State(dev *hid.Device) (Button, error) {
	buf := make([]byte, 8)
	buf[0] = 0x01
	buf[7] = 0x02

	if _, err := dev.Write(buf); err != nil {
		return Button{}, err
	}

	if _, err := dev.ReadTimeout(buf, 200); err != nil {
		return Button{}, err
	}

	if buf[7] != 0x03 {
		return Button{}, errors.New("bad magic")
	}

	return Button{
		buf[0]&(1<<0) == 0,
		buf[0]&(1<<1) == 0,
	}, nil
}

func Poll(dev *hid.Device) <-chan Button {
	ch := make(chan Button)
	go func() {
		for {
			state, err := State(dev)
			if err != nil {
				panic(err)
			}
			ch <- state
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return ch
}

func Open() (*hid.Device, error) {
	return hid.Open(0x1D34, 0x000D, "")
}
