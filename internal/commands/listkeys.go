package commands

import (
	"fmt"
	"strings"

	"github.com/mchampanis/kronctl/internal/protocol"
)

// ListKeys prints all valid key names grouped by category.
func ListKeys() {
	fmt.Println("Physical key positions (first argument of 'remap'):")
	fmt.Println()
	printMatrixRow(0, "Fn row")
	printMatrixRow(1, "Number row")
	printMatrixRow(2, "QWERTY row")
	printMatrixRow(3, "Home row")
	printMatrixRow(4, "Shift row")
	printMatrixRow(5, "Bottom row")
	fmt.Println()

	fmt.Println("Keycodes (second argument of 'remap'):")
	fmt.Println("  Some names (ESC, F1, etc.) also appear as positions above.")
	fmt.Println("  In 'remap', the first argument is always a position, the second a keycode.")
	fmt.Println()
	for _, g := range protocol.AllKeycodes() {
		name := g.Name
		if len(g.Aliases) > 0 {
			name += " (" + strings.Join(g.Aliases, ", ") + ")"
		}
		fmt.Printf("  %-40s 0x%04X (%d)\n", name, g.Code, g.Code)
	}
}

func printMatrixRow(row uint8, label string) {
	fmt.Printf("  Row %d (%s):", row, label)
	for col := uint8(0); col < protocol.Cols; col++ {
		if name := protocol.PosName(row, col); name != "" {
			fmt.Printf(" %s", name)
		}
	}
	fmt.Println()
}
