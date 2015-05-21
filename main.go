package main

import (
	"fmt"
	"log"
	"time"

	"github.com/GeertJohan/go.hid"
)

type Button struct {
	Button bool
	Lid    bool
}

func GetState(dev *hid.Device) Button {
	buf := make([]byte, 8)
	buf[0] = 0x01
	buf[7] = 0x02

	if _, err := dev.Write(buf); err != nil {
		log.Fatal(err)
	}

	if _, err := dev.ReadTimeout(buf, 200); err != nil {
		log.Fatal(err)
	}

	if buf[7] != 0x03 {
		log.Fatal("bad magic")
	}

	return Button{
		buf[0]&(1<<0) == 0,
		buf[0]&(1<<1) == 0,
	}
}

func PollState(dev *hid.Device) <-chan Button {
	state := make(chan Button)
	go func() {
		for {
			state <- GetState(dev)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return state
}

func main() {
	dev, err := hid.Open(0x1D34, 0x000D, "")
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	state := PollState(dev)

	for {
		fmt.Println(<-state)
	}
}
