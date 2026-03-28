package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mchampanis/kronctl/internal/commands"
	"github.com/mchampanis/kronctl/internal/hid"
	"github.com/mchampanis/kronctl/internal/protocol"
)

var version = "dev"

const defaultLayer = 2 // Win base

func usage() {
	fmt.Fprintf(os.Stderr, "kronctl %s - Keychron Q3 Ultra key remapping tool\n\n", version)
	fmt.Fprintln(os.Stderr, "Usage: kronctl <command> [args...]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  firmware-version                Print keyboard firmware version")
	fmt.Fprintln(os.Stderr, "  remap <key> <target> [-l N]     Remap a key (default layer 2)")
	fmt.Fprintln(os.Stderr, "  dump [-l N]                     Dump key mappings (-l -1 for all)")
	fmt.Fprintln(os.Stderr, "  save <file>                     Save all layers to JSON")
	fmt.Fprintln(os.Stderr, "  load <file>                     Restore mappings from JSON")
	fmt.Fprintln(os.Stderr, "  revert                          Undo a load (restore originals)")
	fmt.Fprintln(os.Stderr, "  list-keys                       Print valid key names")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Layers: 0=Mac base, 1=Mac Fn, 2=Win base, 3=Win Fn")
	fmt.Fprintln(os.Stderr)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	// Commands that don't need HID
	switch os.Args[1] {
	case "list-keys":
		commands.ListKeys()
		return
	case "version", "-v", "--version":
		fmt.Println(version)
		return
	case "help", "-h", "--help":
		usage()
		return
	}

	if err := hid.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	defer hid.Exit()

	args := os.Args[2:]

	var err error
	switch os.Args[1] {
	case "firmware-version":
		err = commands.Version()

	case "remap":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: kronctl remap <key> <target> [-l layer]")
			os.Exit(1)
		}
		layer, lerr := parseLayer(args[2:], defaultLayer)
		if lerr != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", lerr)
			os.Exit(1)
		}
		if layer < 0 {
			fmt.Fprintln(os.Stderr, "error: remap requires a specific layer (0-3), not -1")
			os.Exit(1)
		}
		err = commands.Remap(args[0], args[1], layer)

	case "dump":
		layer, lerr := parseLayer(args, defaultLayer)
		if lerr != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", lerr)
			os.Exit(1)
		}
		err = commands.Dump(layer)

	case "save":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "usage: kronctl save <file>")
			os.Exit(1)
		}
		err = commands.Save(args[0])

	case "load":
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "usage: kronctl load <file>")
			os.Exit(1)
		}
		err = commands.Load(args[0])

	case "revert":
		err = commands.Revert()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// parseLayer extracts -l N from args, returning def if not found.
func parseLayer(args []string, def int) (int, error) {
	for i, a := range args {
		if a == "-l" {
			if i+1 >= len(args) {
				return 0, fmt.Errorf("-l requires a layer number")
			}
			n, err := strconv.Atoi(args[i+1])
			if err != nil {
				return 0, fmt.Errorf("invalid layer number: %s", args[i+1])
			}
			if n < -1 || n >= protocol.MaxLayers {
				return 0, fmt.Errorf("layer %d out of range (-1 for all, 0-%d)", n, protocol.MaxLayers-1)
			}
			return n, nil
		}
	}
	return def, nil
}
