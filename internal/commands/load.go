package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mchampanis/kronctl/internal/hid"
	"github.com/mchampanis/kronctl/internal/protocol"
)

const configBase = ".config"
const configDir = "kronctl"
const revertFile = "active.revert.json"

// resolveLayout resolves named keys (key/target fields) into raw row/col/keycode.
// Entries that already use row/col/keycode are left unchanged.
// Mixing named and raw fields for the same axis is rejected.
func resolveLayout(layout *Layout) error {
	for i, ll := range layout.Layers {
		for j, km := range ll.Keys {
			if km.Key != "" && (km.Row != 0 || km.Col != 0) {
				return fmt.Errorf("layer %d: key %q conflicts with row/col -- use one or the other", ll.Layer, km.Key)
			}
			if km.Target != "" && km.Keycode != 0 {
				return fmt.Errorf("layer %d: target %q conflicts with keycode -- use one or the other", ll.Layer, km.Target)
			}
			if km.Key != "" {
				pos, err := protocol.LookupPos(km.Key)
				if err != nil {
					return fmt.Errorf("layer %d, key %q: %w", ll.Layer, km.Key, err)
				}
				layout.Layers[i].Keys[j].Row = pos.Row
				layout.Layers[i].Keys[j].Col = pos.Col
			}
			if km.Target != "" {
				code, err := protocol.LookupKeycode(km.Target)
				if err != nil {
					return fmt.Errorf("layer %d, target %q: %w", ll.Layer, km.Target, err)
				}
				layout.Layers[i].Keys[j].Keycode = code
			}
		}
	}
	return nil
}

// activeRevertPath returns the path to the single revert file in ~/<configBase>/<configDir>/.
func activeRevertPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home directory: %w", err)
	}
	dir := filepath.Join(home, configBase, configDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create config dir: %w", err)
	}
	return filepath.Join(dir, revertFile), nil
}

// Load reads a JSON layout file and writes it to the keyboard.
// Before applying, it captures the current values of affected keys into
// ~/<configBase>/<configDir>/active.revert.json. If a revert file already exists
// (a preset is active), it reverts first, then captures fresh originals
// and applies the new preset.
func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var layout Layout
	if err := json.Unmarshal(data, &layout); err != nil {
		return fmt.Errorf("parse layout: %w", err)
	}

	for _, ll := range layout.Layers {
		if ll.Layer < 0 || ll.Layer >= protocol.MaxLayers {
			return fmt.Errorf("invalid layer %d in layout file (0-%d)", ll.Layer, protocol.MaxLayers-1)
		}
	}

	if err := resolveLayout(&layout); err != nil {
		return err
	}

	dev, err := hid.Open()
	if err != nil {
		return err
	}
	defer dev.Close()

	rp, err := activeRevertPath()
	if err != nil {
		return err
	}

	// If a previous preset is active, revert it first
	if _, err := os.Stat(rp); err == nil {
		if err := applyLayout(dev, rp); err != nil {
			return fmt.Errorf("revert previous preset: %w", err)
		}
		if err := os.Remove(rp); err != nil {
			return fmt.Errorf("remove stale revert file: %w", err)
		}
		fmt.Println("reverted previous preset")
	}

	// Capture current values of keys we're about to overwrite
	originals, err := captureOriginals(dev, layout)
	if err != nil {
		return err
	}
	if len(originals.Layers) > 0 {
		odata, err := json.MarshalIndent(originals, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal revert data: %w", err)
		}
		if err := os.WriteFile(rp, odata, 0644); err != nil {
			return fmt.Errorf("write revert file: %w", err)
		}
	}

	// Apply the preset
	total, err := writeLayout(dev, layout)
	if err != nil {
		return err
	}

	fmt.Printf("loaded %d keys across %d layers from %s\n", total, len(layout.Layers), path)
	return nil
}

// Revert restores original key values from the active revert file.
func Revert() error {
	rp, err := activeRevertPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(rp); os.IsNotExist(err) {
		return fmt.Errorf("no active preset to revert")
	}

	dev, err := hid.Open()
	if err != nil {
		return err
	}
	defer dev.Close()

	if err := applyLayout(dev, rp); err != nil {
		return err
	}
	if err := os.Remove(rp); err != nil {
		return fmt.Errorf("remove revert file: %w", err)
	}
	fmt.Println("reverted to original mappings")
	return nil
}

// captureOriginals reads the current keycodes for every position in the layout.
func captureOriginals(dev *hid.Device, layout Layout) (Layout, error) {
	var originals Layout
	for _, ll := range layout.Layers {
		ol := LayerLayout{Layer: ll.Layer}
		for _, km := range ll.Keys {
			code, err := protocol.GetKeycode(dev, uint8(ll.Layer), km.Row, km.Col)
			if err != nil {
				return Layout{}, fmt.Errorf("read original (%d,%d,%d): %w", ll.Layer, km.Row, km.Col, err)
			}
			ol.Keys = append(ol.Keys, KeyMapping{
				Row:     km.Row,
				Col:     km.Col,
				Keycode: code,
				Name:    protocol.KeycodeName(code),
			})
		}
		originals.Layers = append(originals.Layers, ol)
	}
	return originals, nil
}

// writeLayout sends all keycodes in a Layout to the keyboard.
func writeLayout(dev *hid.Device, layout Layout) (int, error) {
	total := 0
	for _, ll := range layout.Layers {
		for _, km := range ll.Keys {
			if err := protocol.SetKeycode(dev, uint8(ll.Layer), km.Row, km.Col, km.Keycode); err != nil {
				return total, fmt.Errorf("write (%d,%d,%d) (after %d keys): %w", ll.Layer, km.Row, km.Col, total, err)
			}
			total++
		}
	}
	return total, nil
}

// applyLayout reads a layout JSON file and writes it to the keyboard.
func applyLayout(dev *hid.Device, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	var layout Layout
	if err := json.Unmarshal(data, &layout); err != nil {
		return fmt.Errorf("parse layout: %w", err)
	}
	_, err = writeLayout(dev, layout)
	return err
}
