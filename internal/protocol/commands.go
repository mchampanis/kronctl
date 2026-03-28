package protocol

import (
	"fmt"
	"strings"

	"github.com/mchampanis/kronctl/internal/hid"
)

// VIA protocol commands
const (
	CmdGetKeycode    = 0x04
	CmdSetKeycode    = 0x05
	CmdGetLayerCount = 0x11
	CmdFirmwareVer   = 0xA1

	MaxLayers = 4
)

// ValidatePos checks that layer, row, col are within the Q3 Ultra matrix bounds.
func ValidatePos(layer, row, col uint8) error {
	if layer >= MaxLayers {
		return fmt.Errorf("layer %d out of range (0-%d)", layer, MaxLayers-1)
	}
	if row >= Rows {
		return fmt.Errorf("row %d out of range (0-%d)", row, Rows-1)
	}
	if col >= Cols {
		return fmt.Errorf("col %d out of range (0-%d)", col, Cols-1)
	}
	return nil
}

// GetKeycode reads the keycode at (layer, row, col) from the keyboard.
func GetKeycode(dev *hid.Device, layer, row, col uint8) (uint16, error) {
	if err := ValidatePos(layer, row, col); err != nil {
		return 0, err
	}

	var req [hid.ReportLen]byte
	req[0] = CmdGetKeycode
	req[1] = layer
	req[2] = row
	req[3] = col

	resp, err := dev.Send(req)
	if err != nil {
		return 0, err
	}

	if resp[0] != CmdGetKeycode {
		return 0, hid.ErrUnexpectedResponse(resp[0], CmdGetKeycode)
	}
	if resp[1] != layer || resp[2] != row || resp[3] != col {
		return 0, fmt.Errorf("get keycode echo mismatch: sent (%d,%d,%d), got (%d,%d,%d)",
			layer, row, col, resp[1], resp[2], resp[3])
	}

	// Response: [0x04, layer, row, col, keycode_hi, keycode_lo]
	code := uint16(resp[4])<<8 | uint16(resp[5])
	return code, nil
}

// SetKeycode writes a keycode at (layer, row, col) on the keyboard.
func SetKeycode(dev *hid.Device, layer, row, col uint8, keycode uint16) error {
	if err := ValidatePos(layer, row, col); err != nil {
		return err
	}

	var req [hid.ReportLen]byte
	req[0] = CmdSetKeycode
	req[1] = layer
	req[2] = row
	req[3] = col
	req[4] = uint8(keycode >> 8)   // high byte
	req[5] = uint8(keycode & 0xFF) // low byte

	resp, err := dev.Send(req)
	if err != nil {
		return err
	}

	if resp[0] != CmdSetKeycode {
		return hid.ErrUnexpectedResponse(resp[0], CmdSetKeycode)
	}
	// Verify the echo matches our request
	if resp[1] != layer || resp[2] != row || resp[3] != col {
		return fmt.Errorf("set keycode echo mismatch: sent (%d,%d,%d), got (%d,%d,%d)",
			layer, row, col, resp[1], resp[2], resp[3])
	}
	echoCode := uint16(resp[4])<<8 | uint16(resp[5])
	if echoCode != keycode {
		return fmt.Errorf("set keycode echo mismatch: sent 0x%04X, got 0x%04X", keycode, echoCode)
	}
	return nil
}

// GetLayerCount returns the number of layers the keyboard supports.
func GetLayerCount(dev *hid.Device) (int, error) {
	var req [hid.ReportLen]byte
	req[0] = CmdGetLayerCount

	resp, err := dev.Send(req)
	if err != nil {
		return 0, err
	}

	if resp[0] != CmdGetLayerCount {
		return 0, hid.ErrUnexpectedResponse(resp[0], CmdGetLayerCount)
	}

	count := int(resp[1])
	if count < 1 || count > MaxLayers {
		return 0, fmt.Errorf("invalid layer count %d (expected 1-%d)", count, MaxLayers)
	}
	return count, nil
}

// GetFirmwareVersion returns the keyboard's firmware version string.
func GetFirmwareVersion(dev *hid.Device) (string, error) {
	var req [hid.ReportLen]byte
	req[0] = CmdFirmwareVer

	resp, err := dev.Send(req)
	if err != nil {
		return "", fmt.Errorf("version query: %w", err)
	}

	if resp[0] != CmdFirmwareVer {
		return "", hid.ErrUnexpectedResponse(resp[0], CmdFirmwareVer)
	}

	// Firmware string starts at byte 1, null-terminated ASCII
	raw := resp[1:]
	end := len(raw)
	for i, b := range raw {
		if b == 0 {
			end = i
			break
		}
	}
	return strings.TrimSpace(string(raw[:end])), nil
}
