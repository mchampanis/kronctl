# kronctl

CLI tool for remapping Keychron Q3 Ultra keys over HID.

Talks the same protocol as [Keychron Launcher](https://launcher.keychron.com),
but from the command line.

## Install

Requires Go 1.24+ and CGO (MinGW/gcc on PATH for Windows).

```
make release
```

## Usage

```
kronctl firmware-version     # print keyboard firmware version
kronctl remap ESC F1         # remap ESC to F1 on the default layer
kronctl remap ESC F1 -l 3    # remap on layer 3 (Win Fn)
kronctl dump                 # dump current key mappings
kronctl dump -l 2            # dump layer 2 (Win base)
kronctl save layout.json     # save current mappings to file
kronctl load layout.json     # restore mappings from file
kronctl revert               # undo a load (restore originals)
kronctl list-keys            # print all valid key names
```

## Presets

A layout file can contain the full keyboard or just a handful of keys.
`load` only writes the keys present in the file, so you can create small
game-specific presets without touching the rest of the keyboard.

Before applying, `load` captures the original values of affected keys to
`~/.config/kronctl/active.revert.json`. Running `revert` restores them.
If you load a second preset without reverting, the first is automatically
reverted before the new one is applied.

```
kronctl load fps.json        # captures originals, applies preset
kronctl load moba.json       # auto-reverts fps, applies moba
kronctl revert               # back to stock
```

Example preset ([examples/fps.json](examples/fps.json)) -- remaps
CAPSLOCK to LCTRL and ESC to GRAVE on layer 2:

```json
{
  "layers": [
    {
      "layer": 2,
      "keys": [
        {"key": "CAPS", "target": "LCTL"},
        {"key": "ESC", "target": "GRV"}
      ]
    }
  ]
}
```

Hand-written presets use key names (`key` and `target` fields).
Files produced by `kronctl save` use raw matrix positions (`row`, `col`,
`keycode`) -- this is lossless since not every matrix position has a
named key in the map. Both formats are accepted by `load`.

Use `kronctl list-keys` to see valid key and keycode names.

## Currently supported/tested hardware

- Keychron Q3 Ultra 8K ANSI (vendor 0x3434, product 0x1230)

## License

MIT
