package protocol

import (
	"fmt"
	"strings"
)

// Q3 Ultra ANSI matrix dimensions
const (
	Rows = 6
	Cols = 17
)

// Pos is a (row, col) position in the keyboard matrix.
type Pos struct {
	Row uint8
	Col uint8
}

// Q3 Ultra ANSI physical layout.
// Maps physical key labels to matrix positions.
// Confirmed positions from HID captures are marked.
var keyPosByName = map[string]Pos{
	// Row 0: function row + nav
	"ESC": {0, 0}, // confirmed
	"F1":  {0, 1}, // confirmed
	"F2":  {0, 2},
	"F3":  {0, 3},
	"F4":  {0, 4},
	"F5":  {0, 5},
	"F6":  {0, 6},
	"F7":  {0, 7},
	"F8":  {0, 8},
	"F9":  {0, 9},
	"F10": {0, 10},
	"F11": {0, 11},
	"F12": {0, 12},
	"KNOB": {0, 13}, // rotary encoder / mute
	"PSCR": {0, 14},
	"CALC": {0, 15},
	"EMOJ": {0, 16}, // confirmed: RGB_TOG at (0,16) on Fn layer

	// Row 1: number row
	"GRV":  {1, 0},
	"N1":   {1, 1},
	"N2":   {1, 2},
	"N3":   {1, 3},
	"N4":   {1, 4},
	"N5":   {1, 5},
	"N6":   {1, 6},
	"N7":   {1, 7},
	"N8":   {1, 8},
	"N9":   {1, 9},
	"N0":   {1, 10},
	"MINS": {1, 11},
	"EQL":  {1, 12},
	"BSPC": {1, 13},
	"INS":  {1, 14},
	"HOME": {1, 15},
	"PGUP": {1, 16},

	// Row 2: QWERTY row
	"TAB":  {2, 0},
	"Q":    {2, 1},
	"W":    {2, 2},
	"E":    {2, 3},
	"R":    {2, 4},
	"T":    {2, 5},
	"Y":    {2, 6}, // confirmed
	"U":    {2, 7},
	"I":    {2, 8},
	"O":    {2, 9},
	"P":    {2, 10},
	"LBRC": {2, 11},
	"RBRC": {2, 12},
	"BSLS": {2, 13},
	"DEL":  {2, 14},
	"END":  {2, 15},
	"PGDN": {2, 16},

	// Row 3: home row (enter is at col 13, col 12 is empty)
	"CAPS": {3, 0},
	"A":    {3, 1}, // confirmed
	"S":    {3, 2},
	"D":    {3, 3},
	"F":    {3, 4},
	"G":    {3, 5},
	"H":    {3, 6},
	"J":    {3, 7},
	"K":    {3, 8},
	"L":    {3, 9},
	"SCLN": {3, 10},
	"QUOT": {3, 11},
	"ENT":  {3, 13}, // verified by dump

	// Row 4: shift row (col 1 empty, Z starts at col 2)
	"LSFT": {4, 0},
	"Z":    {4, 2},  // verified: dump showed Z keycode at (4,2)
	"X":    {4, 3},
	"C":    {4, 4},
	"V":    {4, 5},
	"B":    {4, 6},
	"N":    {4, 7},
	"M":    {4, 8},
	"COMM": {4, 9},
	"DOT":  {4, 10},
	"SLSH": {4, 11},
	"RSFT": {4, 13}, // verified by dump
	"UP":   {4, 15},

	// Row 5: bottom row (verified by dump)
	"LCTL": {5, 0}, // confirmed
	"LWIN": {5, 1},
	"LALT": {5, 2},
	"SPC":  {5, 6},  // verified by dump
	"RALT": {5, 10},
	"RWIN": {5, 11},
	"MO3":  {5, 12}, // reads as 0x5223 (MO(3) layer tap)
	"RCTL": {5, 13},
	"LEFT": {5, 14},
	"DOWN": {5, 15},
	"RGHT": {5, 16},
}

// LookupPos returns the matrix position for a physical key name (case-insensitive).
func LookupPos(name string) (Pos, error) {
	n := strings.ToUpper(name)
	if pos, ok := keyPosByName[n]; ok {
		return pos, nil
	}
	return Pos{}, fmt.Errorf("unknown key position: %s (use 'kronctl list-keys' to see valid names)", name)
}

// nameByPos is the reverse mapping from matrix position to key name.
var nameByPos map[Pos]string

func init() {
	nameByPos = make(map[Pos]string, len(keyPosByName))
	for name, pos := range keyPosByName {
		existing, ok := nameByPos[pos]
		if !ok || len(name) < len(existing) || (len(name) == len(existing) && name < existing) {
			nameByPos[pos] = name
		}
	}
}

// PosName returns the physical key name for a matrix position, or empty string.
func PosName(row, col uint8) string {
	return nameByPos[Pos{row, col}]
}
