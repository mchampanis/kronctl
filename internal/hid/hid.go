package hid

import (
	"errors"
	"fmt"
	"time"

	gohid "github.com/sstallion/go-hid"
)

// errFound is a sentinel used to break out of HID enumeration early.
var errFound = errors.New("found")

const (
	VendorID  = 0x3434
	ProductID = 0x1230
	ReportLen = 32
	UsagePage = 0xFF60 // QMK/VIA vendor-specific usage page
)

// Device wraps a HID connection to the keyboard.
type Device struct {
	dev *gohid.Device
}

// Init initializes the HID subsystem. Call once at program startup.
func Init() error {
	if err := gohid.Init(); err != nil {
		return fmt.Errorf("hid init: %w", err)
	}
	return nil
}

// Exit tears down the HID subsystem. Call once at program shutdown.
func Exit() {
	gohid.Exit()
}

// Open finds and opens the Keychron Q3 Ultra's vendor-specific HID interface.
// Init must be called before Open.
func Open() (*Device, error) {
	var path string
	err := gohid.Enumerate(VendorID, ProductID, func(info *gohid.DeviceInfo) error {
		if info.UsagePage == UsagePage {
			path = info.Path
			return errFound
		}
		return nil
	})
	if err != nil && !errors.Is(err, errFound) {
		return nil, fmt.Errorf("hid enumerate: %w", err)
	}
	if path == "" {
		return nil, fmt.Errorf("keyboard not found (vendor=%04x product=%04x usage_page=%04x)", VendorID, ProductID, UsagePage)
	}

	dev, err := gohid.OpenPath(path)
	if err != nil {
		return nil, fmt.Errorf("hid open: %w", err)
	}

	return &Device{dev: dev}, nil
}

// Close releases the HID device.
func (d *Device) Close() error {
	return d.dev.Close()
}

// Send writes a 32-byte report and reads the response.
func (d *Device) Send(data [ReportLen]byte) ([ReportLen]byte, error) {
	// Prepend report ID 0x00 for Windows HID
	buf := make([]byte, ReportLen+1)
	buf[0] = 0x00
	copy(buf[1:], data[:])

	if _, err := d.dev.Write(buf); err != nil {
		return [ReportLen]byte{}, fmt.Errorf("hid write: %w", err)
	}

	var resp [ReportLen]byte
	n, err := d.dev.ReadWithTimeout(resp[:], time.Second)
	if err != nil {
		return [ReportLen]byte{}, fmt.Errorf("hid read: %w", err)
	}
	if n == 0 {
		return [ReportLen]byte{}, fmt.Errorf("hid read: timeout")
	}
	if n != ReportLen {
		return [ReportLen]byte{}, fmt.Errorf("hid read: expected %d bytes, got %d", ReportLen, n)
	}

	return resp, nil
}

// ErrUnexpectedResponse returns an error for an unexpected response byte.
func ErrUnexpectedResponse(got, want byte) error {
	return fmt.Errorf("unexpected response: got 0x%02x, want 0x%02x", got, want)
}
