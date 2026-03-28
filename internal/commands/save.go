package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mchampanis/kronctl/internal/hid"
	"github.com/mchampanis/kronctl/internal/protocol"
)

// Layout represents a saved keyboard layout.
type Layout struct {
	Layers []LayerLayout `json:"layers"`
}

// LayerLayout is a single layer's key mappings.
type LayerLayout struct {
	Layer int          `json:"layer"`
	Keys  []KeyMapping `json:"keys"`
}

// KeyMapping is one key's position and assigned keycode.
// For presets, Key and Target can be used instead of row/col/keycode.
type KeyMapping struct {
	Row     uint8  `json:"row"`
	Col     uint8  `json:"col"`
	Keycode uint16 `json:"keycode"`
	Name    string `json:"name,omitempty"`
	Key     string `json:"key,omitempty"`    // physical key name (e.g. "CAPS")
	Target  string `json:"target,omitempty"` // keycode name (e.g. "LCTL")
}

// Save reads the entire keymap and writes it to a JSON file.
func Save(path string) error {
	dev, err := hid.Open()
	if err != nil {
		return err
	}
	defer dev.Close()

	count, err := protocol.GetLayerCount(dev)
	if err != nil {
		return fmt.Errorf("get layer count: %w", err)
	}

	var layout Layout
	total := 0

	for l := range count {
		ll := LayerLayout{Layer: l}
		for row := uint8(0); row < protocol.Rows; row++ {
			for col := uint8(0); col < protocol.Cols; col++ {
				code, err := protocol.GetKeycode(dev, uint8(l), row, col)
				if err != nil {
					return fmt.Errorf("read (%d,%d,%d): %w", l, row, col, err)
				}
				// Skip empty matrix slots (no physical key, no keycode)
				if code == 0 && protocol.PosName(row, col) == "" {
					continue
				}
				ll.Keys = append(ll.Keys, KeyMapping{
					Row:     row,
					Col:     col,
					Keycode: code,
					Name:    protocol.KeycodeName(code),
				})
				total++
			}
		}
		layout.Layers = append(layout.Layers, ll)
	}

	data, err := json.MarshalIndent(layout, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	fmt.Printf("saved %d keys across %d layers to %s\n", total, count, path)
	return nil
}
