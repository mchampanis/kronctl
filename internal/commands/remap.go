package commands

import (
	"fmt"

	"github.com/mchampanis/kronctl/internal/hid"
	"github.com/mchampanis/kronctl/internal/protocol"
)

// Remap sets a physical key to produce a different keycode.
func Remap(keyName, targetName string, layer int) error {
	pos, err := protocol.LookupPos(keyName)
	if err != nil {
		return err
	}

	keycode, err := protocol.LookupKeycode(targetName)
	if err != nil {
		return err
	}

	dev, err := hid.Open()
	if err != nil {
		return err
	}
	defer dev.Close()

	if err := protocol.SetKeycode(dev, uint8(layer), pos.Row, pos.Col, keycode); err != nil {
		return fmt.Errorf("remap failed: %w", err)
	}

	fmt.Printf("%s -> %s (layer %d, pos %d,%d, keycode 0x%04X)\n",
		keyName, protocol.KeycodeName(keycode), layer, pos.Row, pos.Col, keycode)
	return nil
}
