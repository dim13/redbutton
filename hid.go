package redbutton

import (
	"errors"
	"io"

	"github.com/sstallion/go-hid"
)

const (
	vendorID  = 0x1d34
	productID = 0x000d
)

var (
	ErrUnsupported = errors.New("unsupproted platform")
	ErrNotFound    = errors.New("device not found")
)

func Open() (io.ReadWriteCloser, error) {
	if err := hid.Init(); err != nil {
		return nil, err
	}
	return hid.OpenFirst(vendorID, productID)
}
