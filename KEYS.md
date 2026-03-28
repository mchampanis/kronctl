# Key names

## Physical key positions

Used as the first argument to `kronctl remap` to identify which physical key to remap.

```
Row 0 (Fn row):     ESC F1 F2 F3 F4 F5 F6 F7 F8 F9 F10 F11 F12 KNOB PSCR CALC EMOJ
Row 1 (Number row): GRV N1 N2 N3 N4 N5 N6 N7 N8 N9 N0 MINS EQL BSPC INS HOME PGUP
Row 2 (QWERTY row): TAB Q W E R T Y U I O P LBRC RBRC BSLS DEL END PGDN
Row 3 (Home row):   CAPS A S D F G H J K L SCLN QUOT ENT
Row 4 (Shift row):  LSFT Z X C V B N M COMM DOT SLSH RSFT UP
Row 5 (Bottom row): LCTL LWIN LALT SPC RALT RWIN MO3 RCTL LEFT DOWN RGHT
```

## Keycodes

Used as the second argument to `kronctl remap` to set what the key produces.
Names are case-insensitive. The `KC_` prefix is optional.

### Basic

| Name | Aliases | Hex |
|------|---------|-----|
| A-Z | | 0x04-0x1D |
| 0-9 | | 0x1E-0x27 |
| ENT | ENTER | 0x28 |
| ESC | ESCAPE | 0x29 |
| BSPC | BACKSPACE | 0x2A |
| TAB | | 0x2B |
| SPC | SPACE | 0x2C |
| MINS | MINUS | 0x2D |
| EQL | EQUAL | 0x2E |
| LBRC | LBRACKET | 0x2F |
| RBRC | RBRACKET | 0x30 |
| BSLS | BACKSLASH | 0x31 |
| SCLN | SEMICOLON | 0x33 |
| QUOT | QUOTE | 0x34 |
| GRV | GRAVE | 0x35 |
| COMM | COMMA | 0x36 |
| DOT | | 0x37 |
| SLSH | SLASH | 0x38 |
| CAPS | CAPSLOCK | 0x39 |
| F1-F12 | | 0x3A-0x45 |
| F13-F24 | | 0x68-0x73 |
| PSCR | PRINTSCREEN | 0x46 |
| SLCK | SCROLLLOCK | 0x47 |
| PAUS | PAUSE | 0x48 |
| INS | INSERT | 0x49 |
| HOME | | 0x4A |
| PGUP | PAGEUP | 0x4B |
| DEL | DELETE | 0x4C |
| END | | 0x4D |
| PGDN | PAGEDOWN | 0x4E |
| RIGHT | RGHT | 0x4F |
| LEFT | | 0x50 |
| DOWN | | 0x51 |
| UP | | 0x52 |
| APP | MENU | 0x65 |

### Modifiers

| Name | Aliases | Hex |
|------|---------|-----|
| LCTL | LCTRL | 0xE0 |
| LSFT | LSHIFT | 0xE1 |
| LALT | | 0xE2 |
| LGUI | LWIN, LCMD | 0xE3 |
| RCTL | RCTRL | 0xE4 |
| RSFT | RSHIFT | 0xE5 |
| RALT | | 0xE6 |
| RGUI | RWIN, RCMD | 0xE7 |

### Media

| Name | Aliases | Hex |
|------|---------|-----|
| MUTE | | 0x7F |
| VOLU | VOLUP | 0xA9 |
| VOLD | VOLDN | 0xAA |
| MNXT | MEDNEXT | 0xAB |
| MPRV | MEDPREV | 0xAC |
| MSTP | MEDSTOP | 0xAD |
| MPLY | MEDPLAY | 0xAE |
| MSEL | | 0xAF |
| EJCT | EJECT | 0xB0 |
| MFFD | MEDFASTFWD | 0xBB |
| MRWD | MEDREWIND | 0xBC |
| BRIU | BRIGHTNESSUP | 0xBD |
| BRID | BRIGHTNESSDN | 0xBE |

### Web / application

| Name | Aliases | Hex |
|------|---------|-----|
| WSRC | WWW_SEARCH | 0xB4 |
| WHOM | WWW_HOME | 0xB5 |
| WBAK | WWW_BACK | 0xB6 |
| WFWD | WWW_FORWARD | 0xB7 |
| WSTP | WWW_STOP | 0xB8 |
| WREF | WWW_REFRESH | 0xB9 |
| WFAV | WWW_FAVORITES | 0xBA |
| MAIL | | 0xB1 |
| CALC | CALCULATOR | 0xB2 |
| MYCM | MYCOMPUTER | 0xB3 |
| MIC | MICROPHONE | 0xC1 |
| LPAD | LAUNCHPAD | 0xC2 |

### Mouse

| Name | Hex |
|------|-----|
| MS_UP, MS_DOWN, MS_LEFT, MS_RIGHT | 0xF0-0xF3 |
| MS_BTN1-MS_BTN5 | 0xF4-0xF8 |
| MS_WH_UP, MS_WH_DOWN | 0xF9-0xFA |
| MS_WH_LEFT, MS_WH_RIGHT | 0xFB-0xFC |

### RGB / backlight

| Name | Description | Hex |
|------|-------------|-----|
| RGB_TOG | Toggle RGB | 0x7820 |
| RGB_MOD | Mode + | 0x7821 |
| RGB_RMOD | Mode - | 0x7822 |
| RGB_HUI | Hue + | 0x7823 |
| RGB_HUD | Hue - | 0x7824 |
| RGB_SAI | Saturation + | 0x7825 |
| RGB_SAD | Saturation - | 0x7826 |
| RGB_VAI | Brightness + | 0x7827 |
| RGB_VAD | Brightness - | 0x7828 |
| RGB_SPI | Effect speed + | 0x7829 |
| RGB_SPD | Effect speed - | 0x782A |
| BL_ON | Backlight on | 0x7800 |
| BL_OFF | Backlight off | 0x7801 |
| BL_TOGG | Backlight toggle | 0x7802 |
| BL_DEC | Backlight - | 0x7803 |
| BL_INC | Backlight + | 0x7804 |
| BL_STEP | Backlight step | 0x7805 |

### Layer functions

| Name | Description | Hex |
|------|-------------|-----|
| MO0-MO3 | Momentary layer (hold) | 0x5220-0x5223 |
| TG0-TG3 | Toggle layer | 0x5260-0x5263 |
| TO0-TO3 | Switch to layer | 0x5210-0x5213 |
| DF0-DF3 | Set default layer | 0x5240-0x5243 |
| OSL0-OSL3 | One-shot layer | 0x5280-0x5283 |
| TT0-TT3 | Layer tap toggle | 0x52C0-0x52C3 |
| FN_MO13 | Fn combo (layers 1+3) | 0x7700 |
| FN_MO23 | Fn combo (layers 2+3) | 0x7701 |

### Keychron custom (0x7E00+)

| Name | Description | Hex |
|------|-------------|-----|
| BT1 | Bluetooth profile 1 | 0x7E00 |
| BT2 | Bluetooth profile 2 | 0x7E01 |
| BT3 | Bluetooth profile 3 | 0x7E02 |
| BOOT | Bootloader mode | 0x7E03 |
| 2_4G | 2.4G wireless | 0x7E04 |
| SCRLK_KC | Scroll Lock (Keychron) | 0x7E05 |
| CMD_COMM | Cmd-Comma (Mac Prefs) | 0x7E06 |
| KC_LOPT | macOS Left Option | 0x7E07 |
| KC_ROPT | macOS Right Option | 0x7E08 |
| KC_LCMD | macOS Left Command | 0x7E09 |
| KC_RCMD | macOS Right Command | 0x7E0A |
| DSKL_M | Desktop Left (Mac) | 0x7E0B |
| DSKR_M | Desktop Right (Mac) | 0x7E0C |
| EMOJI_M | Emoji (Mac) | 0x7E0D |
| TASK | Task View (Win) | 0x7E0E |
| DSKL_W | Desktop Left (Win) | 0x7E0F |
| DSKR_W | Desktop Right (Win) | 0x7E10 |
| FILE | File Explorer (Win) | 0x7E11 |
| WINLOCK | Lock Windows key | 0x7E12 |
| SET | Settings | 0x7E13 |
| EMOJI_W | Emoji (Win) | 0x7E14 |
| PSCR_W | Print Screen (Win) | 0x7E15 |
| SSHOT | Screenshot (Mac) | 0x7E16 |
| BATT | Battery status | 0x7E17 |
| SIRI | Siri / Dictation | 0x7E18 |
| CORTANA | Cortana | 0x7E19 |
| MCTL | Mission Control (Mac) | 0x7E1A |
| KC_LPAD | Launchpad (Mac) | 0x7E1B |
| KC_MUTE_KNOB | Rotary encoder mute | 0x00A8 |

### System

| Name | Description | Hex |
|------|-------------|-----|
| NO | No key (disabled) | 0x0000 |
| TRNS | Transparent (inherit) | 0x0001 |
| NKRO | Toggle N-key rollover | 0x7013 |
| RESET | Reset keyboard | 0x7C00 |
| QK_CLEAR_EEPROM | Clear EEPROM | 0x7C03 |
