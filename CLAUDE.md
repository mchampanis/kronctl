# kronctl

CLI tool for remapping Keychron Q3 Ultra keys over HID without reflashing firmware.

## Build

Requires Go 1.24+ and CGO (MinGW/gcc on PATH for Windows).

```
make          # debug build with race detector
make release  # stripped release build
```

## Architecture

```
cmd/kronctl/         CLI entry point, arg dispatch
internal/hid/        HID device open/close, report send/receive
internal/commands/   One file per command category
```

Follows the same patterns as wincon: no CLI framework, manual arg parsing,
switch-based dispatch, one function per command.

## Protocol

Communicates with the keyboard over USB HID (vendor-specific usage page 0xFF60).
Reports are 32 bytes. The protocol is VIA-compatible with Keychron extensions.

Key commands:
- 0x04: get keycode (layer, row, col) -> keycode
- 0x05: set keycode (layer, row, col, keycode_hi, keycode_lo)
- 0x11: get layer count
- 0xA1: get firmware version string
- 0xA3: keepalive / poll

Keycodes are 16-bit, sent big-endian in bytes 4-5 of the remap packet.
Standard USB HID keycodes fit in the low byte; extended codes (RGB, macros)
use the full 16 bits.

## Device

- Vendor ID: 0x3434
- Product ID: 0x1230
- Name: Keychron Q3 Ultra 8K ANSI
- Layers: 0=Mac base, 1=Mac Fn, 2=Win base, 3=Win Fn
