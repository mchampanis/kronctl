package commands

import (
	"fmt"

	"github.com/mchampanis/kronctl/internal/hid"
	"github.com/mchampanis/kronctl/internal/protocol"
)

// Version queries the keyboard firmware version string.
func Version() error {
	dev, err := hid.Open()
	if err != nil {
		return err
	}
	defer dev.Close()

	version, err := protocol.GetFirmwareVersion(dev)
	if err != nil {
		return err
	}

	fmt.Println(version)
	return nil
}
