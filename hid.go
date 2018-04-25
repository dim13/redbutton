package redbutton

import (
	"errors"

	"github.com/karalabe/hid"
)

const (
	vendorID  = 0x1d34
	productID = 0x000d
)

var (
	ErrUnsupported = errors.New("unsupproted platform")
	ErrNotFound    = errors.New("device not found")
)

func Open() (*hid.Device, error) {
	if !hid.Supported() {
		return nil, ErrUnsupported
	}
	for _, dev := range hid.Enumerate(vendorID, productID) {
		return dev.Open()
	}
	return nil, ErrNotFound
}
