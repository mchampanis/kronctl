package protocol

import (
	"fmt"
	"sort"
	"strings"
)

// Keycode values match the Keychron/QMK protocol encoding (wire format).
// Basic keys (0-231) follow USB HID usage tables.
// Extended keys (>255) are 16-bit Keychron/QMK quantum keycodes.
var keycodeByName = map[string]uint16{
	// No key / transparent
	"NO":   0,
	"TRNS": 1, // transparent (inherit from lower layer)

	// Letters
	"A": 4, "B": 5, "C": 6, "D": 7, "E": 8, "F": 9, "G": 10,
	"H": 11, "I": 12, "J": 13, "K": 14, "L": 15, "M": 16, "N": 17,
	"O": 18, "P": 19, "Q": 20, "R": 21, "S": 22, "T": 23, "U": 24,
	"V": 25, "W": 26, "X": 27, "Y": 28, "Z": 29,

	// Numbers
	"1": 30, "2": 31, "3": 32, "4": 33, "5": 34,
	"6": 35, "7": 36, "8": 37, "9": 38, "0": 39,

	// Standard keys
	"ENT":  40, "ENTER": 40,
	"ESC":  41, "ESCAPE": 41,
	"BSPC": 42, "BACKSPACE": 42,
	"TAB":  43,
	"SPC":  44, "SPACE": 44,
	"MINS": 45, "MINUS": 45,
	"EQL":  46, "EQUAL": 46,
	"LBRC": 47, "LBRACKET": 47,
	"RBRC": 48, "RBRACKET": 48,
	"BSLS": 49, "BACKSLASH": 49,
	"SCLN": 51, "SEMICOLON": 51,
	"QUOT": 52, "QUOTE": 52,
	"GRV":  53, "GRAVE": 53,
	"COMM": 54, "COMMA": 54,
	"DOT":  55,
	"SLSH": 56, "SLASH": 56,
	"CAPS": 57, "CAPSLOCK": 57,

	// Function keys
	"F1": 58, "F2": 59, "F3": 60, "F4": 61, "F5": 62, "F6": 63,
	"F7": 64, "F8": 65, "F9": 66, "F10": 67, "F11": 68, "F12": 69,
	"F13": 104, "F14": 105, "F15": 106, "F16": 107,
	"F17": 108, "F18": 109, "F19": 110, "F20": 111,
	"F21": 112, "F22": 113, "F23": 114, "F24": 115,

	// Navigation
	"PSCR":  70, "PRINTSCREEN": 70,
	"SLCK":  71, "SCROLLLOCK": 71,
	"PAUS":  72, "PAUSE": 72,
	"INS":   73, "INSERT": 73,
	"HOME":  74,
	"PGUP":  75, "PAGEUP": 75,
	"DEL":   76, "DELETE": 76,
	"END":   77,
	"PGDN":  78, "PAGEDOWN": 78,
	"RIGHT": 79, "RGHT": 79,
	"LEFT":  80,
	"DOWN":  81,
	"UP":    82,

	// Numpad
	"NLCK": 83, "NUMLOCK": 83,
	"PSLS": 84, "KP_SLASH": 84,
	"PAST": 85, "KP_ASTERISK": 85,
	"PMNS": 86, "KP_MINUS": 86,
	"PPLS": 87, "KP_PLUS": 87,
	"PENT": 88, "KP_ENTER": 88,
	"P1":   89, "KP_1": 89,
	"P2":   90, "KP_2": 90,
	"P3":   91, "KP_3": 91,
	"P4":   92, "KP_4": 92,
	"P5":   93, "KP_5": 93,
	"P6":   94, "KP_6": 94,
	"P7":   95, "KP_7": 95,
	"P8":   96, "KP_8": 96,
	"P9":   97, "KP_9": 97,
	"P0":   98, "KP_0": 98,
	"PDOT": 99, "KP_DOT": 99,

	// Modifiers
	"LCTL":  224, "LCTRL": 224,
	"LSFT":  225, "LSHIFT": 225,
	"LALT":  226,
	"LGUI":  227, "LWIN": 227, "LCMD": 227,
	"RCTL":  228, "RCTRL": 228,
	"RSFT":  229, "RSHIFT": 229,
	"RALT":  230,
	"RGUI":  231, "RWIN": 231, "RCMD": 231,

	// Application
	"APP": 101, "MENU": 101,

	// Media / consumer (Keychron protocol values)
	"MUTE": 127,
	"VOLU": 169, "VOLUP": 169,
	"VOLD": 170, "VOLDN": 170,
	"MNXT": 171, "MEDNEXT": 171,
	"MPRV": 172, "MEDPREV": 172,
	"MSTP": 173, "MEDSTOP": 173,
	"MPLY": 174, "MEDPLAY": 174,
	"MSEL": 175,
	"EJCT": 176, "EJECT": 176,
	"MAIL": 177,
	"CALC": 178, "CALCULATOR": 178,
	"MYCM": 179, "MYCOMPUTER": 179,
	"MFFD": 187, "MEDFASTFWD": 187,
	"MRWD": 188, "MEDREWIND": 188,
	"BRIU": 189, "BRIGHTNESSUP": 189,
	"BRID": 190, "BRIGHTNESSDN": 190,

	// Web
	"WSRC": 180, "WWW_SEARCH": 180,
	"WHOM": 181, "WWW_HOME": 181,
	"WBAK": 182, "WWW_BACK": 182,
	"WFWD": 183, "WWW_FORWARD": 183,
	"WSTP": 184, "WWW_STOP": 184,
	"WREF": 185, "WWW_REFRESH": 185,
	"WFAV": 186, "WWW_FAVORITES": 186,

	// System
	"MIC":  193, "MICROPHONE": 193,
	"LPAD": 194, "LAUNCHPAD": 194,

	// Mouse keys
	"MS_UP":    240, "MS_DOWN": 241, "MS_LEFT": 242, "MS_RIGHT": 243,
	"MS_BTN1":  244, "MS_BTN2": 245, "MS_BTN3": 246,
	"MS_BTN4":  247, "MS_BTN5": 248,
	"MS_WH_UP": 249, "MS_WH_DOWN": 250,
	"MS_WH_LEFT": 251, "MS_WH_RIGHT": 252,

	// RGB (16-bit extended keycodes)
	"RGB_TOG":  30752,
	"RGB_MOD":  30753,
	"RGB_RMOD": 30754,
	"RGB_HUI":  30755,
	"RGB_HUD":  30756,
	"RGB_SAI":  30757,
	"RGB_SAD":  30758,
	"RGB_VAI":  30759,
	"RGB_VAD":  30760,
	"RGB_SPI":  30761,
	"RGB_SPD":  30762,

	// Backlight
	"BL_ON":   30720,
	"BL_OFF":  30721,
	"BL_TOGG": 30722,
	"BL_DEC":  30723,
	"BL_INC":  30724,
	"BL_STEP": 30725,

	// Keychron-specific
	"KC_MUTE_KNOB": 0x00A8, // rotary encoder mute

	// Keychron QK_KB custom keycodes (0x7E00 + n)
	// Mapped from Custom tab in Launcher, sequential CUSTOM(0)-CUSTOM(27)
	"BT1":      0x7E00, // CUSTOM(0)  Bluetooth profile 1
	"BT2":      0x7E01, // CUSTOM(1)  Bluetooth profile 2
	"BT3":      0x7E02, // CUSTOM(2)  Bluetooth profile 3
	"BOOT":     0x7E03, // CUSTOM(3)  Bootloader mode
	"2_4G":     0x7E04, // CUSTOM(4)  2.4G wireless
	"SCRLK_KC": 0x7E05, // CUSTOM(5)  Scroll Lock (Keychron)
	"CMD_COMM": 0x7E06, // CUSTOM(6)  Cmd-Comma (Mac Preferences)
	"KC_LOPT":  0x7E07, // CUSTOM(7)  macOS Left Option
	"KC_ROPT":  0x7E08, // CUSTOM(8)  macOS Right Option
	"KC_LCMD":  0x7E09, // CUSTOM(9)  macOS Left Command
	"KC_RCMD":  0x7E0A, // CUSTOM(10) macOS Right Command
	"DSKL_M":   0x7E0B, // CUSTOM(11) Desktop Left (Mac)
	"DSKR_M":   0x7E0C, // CUSTOM(12) Desktop Right (Mac)
	"EMOJI_M":  0x7E0D, // CUSTOM(13) Emoji picker (Mac)
	"TASK":     0x7E0E, // CUSTOM(14) Task View (Win+Tab)
	"DSKL_W":   0x7E0F, // CUSTOM(15) Desktop Left (Win)
	"DSKR_W":   0x7E10, // CUSTOM(16) Desktop Right (Win)
	"FILE":     0x7E11, // CUSTOM(17) File Explorer (Win)
	"WINLOCK":  0x7E12, // CUSTOM(18) Lock Windows key
	"SET":      0x7E13, // CUSTOM(19) Settings
	"EMOJI_W":  0x7E14, // CUSTOM(20) Emoji picker (Win)
	"PSCR_W":   0x7E15, // CUSTOM(21) Print Screen (Win)
	"SSHOT":    0x7E16, // CUSTOM(22) Screenshot (Mac)
	"BATT":     0x7E17, // CUSTOM(23) Battery status
	"SIRI":     0x7E18, // CUSTOM(24) Siri / Dictation
	"CORTANA":  0x7E19, // CUSTOM(25) Cortana
	"MCTL":     0x7E1A, // CUSTOM(26) Mission Control (Mac)
	"KC_LPAD":  0x7E1B, // CUSTOM(27) Launchpad (Mac)

	// Layer functions (QMK quantum keycodes)
	"MO0": 0x5220, // MO(0) - momentary layer 0
	"MO1": 0x5221, // MO(1) - momentary layer 1
	"MO2": 0x5222, // MO(2) - momentary layer 2
	"MO3": 0x5223, // MO(3) - momentary layer 3

	"TG0": 0x5260, // TG(0) - toggle layer 0
	"TG1": 0x5261, // TG(1) - toggle layer 1
	"TG2": 0x5262, // TG(2) - toggle layer 2
	"TG3": 0x5263, // TG(3) - toggle layer 3

	"TO0": 0x5210, // TO(0) - switch to layer 0
	"TO1": 0x5211, // TO(1) - switch to layer 1
	"TO2": 0x5212, // TO(2) - switch to layer 2
	"TO3": 0x5213, // TO(3) - switch to layer 3

	"DF0": 0x5240, // DF(0) - set default layer 0
	"DF1": 0x5241, // DF(1) - set default layer 1
	"DF2": 0x5242, // DF(2) - set default layer 2
	"DF3": 0x5243, // DF(3) - set default layer 3

	"OSL0": 0x5280, // OSL(0) - one-shot layer 0
	"OSL1": 0x5281, // OSL(1) - one-shot layer 1
	"OSL2": 0x5282, // OSL(2) - one-shot layer 2
	"OSL3": 0x5283, // OSL(3) - one-shot layer 3

	"TT0": 0x52C0, // TT(0) - layer tap toggle 0
	"TT1": 0x52C1, // TT(1) - layer tap toggle 1
	"TT2": 0x52C2, // TT(2) - layer tap toggle 2
	"TT3": 0x52C3, // TT(3) - layer tap toggle 3

	// Layer Fn combos
	"FN_MO13": 30464,
	"FN_MO23": 30465,

	// NKRO
	"NKRO": 28691, // Toggle N-key rollover

	// System
	"RESET":           31744,
	"QK_CLEAR_EEPROM": 31747,
}

// nameByKeycode is the reverse mapping (keycode -> canonical short name).
var nameByKeycode map[uint16]string

func init() {
	nameByKeycode = make(map[uint16]string)
	// Prefer shorter names; break ties alphabetically for determinism
	for name, code := range keycodeByName {
		existing, ok := nameByKeycode[code]
		if !ok || len(name) < len(existing) || (len(name) == len(existing) && name < existing) {
			nameByKeycode[code] = name
		}
	}
}

// LookupKeycode returns the keycode value for a name (case-insensitive).
// Accepts both short ("ESC") and long ("ESCAPE") forms, with or without "KC_" prefix.
func LookupKeycode(name string) (uint16, error) {
	upper := strings.ToUpper(name)
	// Try exact match first (handles KC_-prefixed entries like KC_LOPT)
	if code, ok := keycodeByName[upper]; ok {
		return code, nil
	}
	// Try with KC_ prefix stripped (handles KC_ESC -> ESC, etc.)
	n := strings.TrimPrefix(upper, "KC_")
	if code, ok := keycodeByName[n]; ok {
		return code, nil
	}
	return 0, fmt.Errorf("unknown keycode: %s", name)
}

// KeycodeName returns the canonical name for a keycode value.
func KeycodeName(code uint16) string {
	if name, ok := nameByKeycode[code]; ok {
		return name
	}
	return fmt.Sprintf("0x%04X", code)
}

// KeycodeGroup is a keycode with its canonical name and any aliases.
type KeycodeGroup struct {
	Name    string   // canonical (shortest) name
	Aliases []string // other names, sorted
	Code    uint16
}

// AllKeycodes returns all keycodes grouped by value, sorted by canonical name.
func AllKeycodes() []KeycodeGroup {
	// Collect aliases per keycode
	byCode := make(map[uint16][]string)
	for name := range keycodeByName {
		code := keycodeByName[name]
		byCode[code] = append(byCode[code], name)
	}

	groups := make([]KeycodeGroup, 0, len(byCode))
	for code, names := range byCode {
		sort.Strings(names)
		canonical := nameByKeycode[code]
		var aliases []string
		for _, n := range names {
			if n != canonical {
				aliases = append(aliases, n)
			}
		}
		groups = append(groups, KeycodeGroup{
			Name:    canonical,
			Aliases: aliases,
			Code:    code,
		})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	return groups
}
