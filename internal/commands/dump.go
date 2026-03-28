package commands

import (
	"fmt"

	"github.com/mchampanis/kronctl/internal/hid"
	"github.com/mchampanis/kronctl/internal/protocol"
)

// Dump reads and prints the keymap for a given layer.
func Dump(layer int) error {
	dev, err := hid.Open()
	if err != nil {
		return err
	}
	defer dev.Close()

	if layer < 0 {
		// Dump all layers
		count, err := protocol.GetLayerCount(dev)
		if err != nil {
			return fmt.Errorf("get layer count: %w", err)
		}
		for l := 0; l < count; l++ {
			if l > 0 {
				fmt.Println()
			}
			if err := dumpLayer(dev, l); err != nil {
				return err
			}
		}
		return nil
	}

	return dumpLayer(dev, layer)
}

func dumpLayer(dev *hid.Device, layer int) error {
	fmt.Printf("Layer %d:\n", layer)

	for row := uint8(0); row < protocol.Rows; row++ {
		for col := uint8(0); col < protocol.Cols; col++ {
			code, err := protocol.GetKeycode(dev, uint8(layer), row, col)
			if err != nil {
				return fmt.Errorf("read (%d,%d,%d): %w", layer, row, col, err)
			}

			if code == 0 {
				continue
			}

			posName := protocol.PosName(row, col)
			codeName := protocol.KeycodeName(code)

			if posName != "" {
				fmt.Printf("  %-6s (%d,%2d) -> %s\n", posName, row, col, codeName)
			} else {
				fmt.Printf("  ?      (%d,%2d) -> %s\n", row, col, codeName)
			}
		}
	}

	return nil
}
