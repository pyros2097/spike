// Copyright 2015 pyros2097. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Keyboard Keys
package keys

type KeyCode uint8

const (
	ANY_KEY             KeyCode = -1
	NUM_0                       = 7
	NUM_1                       = 8
	NUM_2                       = 9
	NUM_3                       = 10
	NUM_4                       = 11
	NUM_5                       = 12
	NUM_6                       = 13
	NUM_7                       = 14
	NUM_8                       = 15
	NUM_9                       = 16
	A                           = 29
	ALT_LEFT                    = 57
	ALT_RIGHT                   = 58
	APOSTROPHE                  = 75
	AT                          = 77
	B                           = 30
	BACK                        = 4
	BACKSLASH                   = 73
	C                           = 31
	CALL                        = 5
	CAMERA                      = 27
	CLEAR                       = 28
	COMMA                       = 55
	D                           = 32
	DEL                         = 67
	BACKSPACE                   = 67
	FORWARD_DEL                 = 112
	DPAD_CENTER                 = 23
	DPAD_DOWN                   = 20
	DPAD_LEFT                   = 21
	DPAD_RIGHT                  = 22
	DPAD_UP                     = 19
	CENTER                      = 23
	DOWN                        = 20
	LEFT                        = 21
	RIGHT                       = 22
	UP                          = 19
	E                           = 33
	ENDCALL                     = 6
	ENTER                       = 66
	ENVELOPE                    = 65
	EQUALS                      = 70
	EXPLORER                    = 64
	F                           = 34
	FOCUS                       = 80
	G                           = 35
	GRAVE                       = 68
	H                           = 36
	HEADSETHOOK                 = 79
	HOME                        = 3
	I                           = 37
	J                           = 38
	K                           = 39
	L                           = 40
	LEFT_BRACKET                = 71
	M                           = 41
	MEDIA_FAST_FORWARD          = 90
	MEDIA_NEXT                  = 87
	MEDIA_PLAY_PAUSE            = 85
	MEDIA_PREVIOUS              = 88
	MEDIA_REWIND                = 89
	MEDIA_STOP                  = 86
	MENU                        = 82
	MINUS                       = 69
	MUTE                        = 91
	N                           = 42
	NOTIFICATION                = 83
	NUM                         = 78
	O                           = 43
	P                           = 44
	PERIOD                      = 56
	PLUS                        = 81
	POUND                       = 18
	POWER                       = 26
	Q                           = 45
	R                           = 46
	RIGHT_BRACKET               = 72
	S                           = 47
	SEARCH                      = 84
	SEMICOLON                   = 74
	SHIFT_LEFT                  = 59
	SHIFT_RIGHT                 = 60
	SLASH                       = 76
	SOFT_LEFT                   = 1
	SOFT_RIGHT                  = 2
	SPACE                       = 62
	STAR                        = 17
	SYM                         = 63
	T                           = 48
	TAB                         = 61
	U                           = 49
	UNKNOWN                     = 0
	V                           = 50
	VOLUME_DOWN                 = 25
	VOLUME_UP                   = 24
	W                           = 51
	X                           = 52
	Y                           = 53
	Z                           = 54
	META_ALT_LEFT_ON            = 16
	META_ALT_ON                 = 2
	META_ALT_RIGHT_ON           = 32
	META_SHIFT_LEFT_ON          = 64
	META_SHIFT_ON               = 1
	META_SHIFT_RIGHT_ON         = 128
	META_SYM_ON                 = 4
	CONTROL_LEFT                = 129
	CONTROL_RIGHT               = 130
	ESCAPE                      = 131
	END                         = 132
	INSERT                      = 133
	PAGE_UP                     = 92
	PAGE_DOWN                   = 93
	PICTSYMBOLS                 = 94
	SWITCH_CHARSET              = 95
	BUTTON_CIRCLE               = 255
	BUTTON_A                    = 96
	BUTTON_B                    = 97
	BUTTON_C                    = 98
	BUTTON_X                    = 99
	BUTTON_Y                    = 100
	BUTTON_Z                    = 101
	BUTTON_L1                   = 102
	BUTTON_R1                   = 103
	BUTTON_L2                   = 104
	BUTTON_R2                   = 105
	BUTTON_THUMBL               = 106
	BUTTON_THUMBR               = 107
	BUTTON_START                = 108
	BUTTON_SELECT               = 109
	BUTTON_MODE                 = 110

	NUMPAD_0 = 144
	NUMPAD_1 = 145
	NUMPAD_2 = 146
	NUMPAD_3 = 147
	NUMPAD_4 = 148
	NUMPAD_5 = 149
	NUMPAD_6 = 150
	NUMPAD_7 = 151
	NUMPAD_8 = 152
	NUMPAD_9 = 153

	COLON = 243
	F1    = 244
	F2    = 245
	F3    = 246
	F4    = 247
	F5    = 248
	F6    = 249
	F7    = 250
	F8    = 251
	F9    = 252
	F10   = 253
	F11   = 254
	F12   = 255
)

// return a human readable representation of the keycode. The returned value can be used in
// {@link Input.Keys#valueOf(String)}
func ToString(keycode uint8) string {
	switch keycode {
	// META* variables should not be used with this method.
	case UNKNOWN:
		return "Unknown"
	case SOFT_LEFT:
		return "Soft Left"
	case SOFT_RIGHT:
		return "Soft Right"
	case HOME:
		return "Home"
	case BACK:
		return "Back"
	case CALL:
		return "Call"
	case ENDCALL:
		return "End Call"
	case NUM_0:
		return "0"
	case NUM_1:
		return "1"
	case NUM_2:
		return "2"
	case NUM_3:
		return "3"
	case NUM_4:
		return "4"
	case NUM_5:
		return "5"
	case NUM_6:
		return "6"
	case NUM_7:
		return "7"
	case NUM_8:
		return "8"
	case NUM_9:
		return "9"
	case STAR:
		return "*"
	case POUND:
		return "#"
	case UP:
		return "Up"
	case DOWN:
		return "Down"
	case LEFT:
		return "Left"
	case RIGHT:
		return "Right"
	case CENTER:
		return "Center"
	case VOLUME_UP:
		return "Volume Up"
	case VOLUME_DOWN:
		return "Volume Down"
	case POWER:
		return "Power"
	case CAMERA:
		return "Camera"
	case CLEAR:
		return "Clear"
	case A:
		return "A"
	case B:
		return "B"
	case C:
		return "C"
	case D:
		return "D"
	case E:
		return "E"
	case F:
		return "F"
	case G:
		return "G"
	case H:
		return "H"
	case I:
		return "I"
	case J:
		return "J"
	case K:
		return "K"
	case L:
		return "L"
	case M:
		return "M"
	case N:
		return "N"
	case O:
		return "O"
	case P:
		return "P"
	case Q:
		return "Q"
	case R:
		return "R"
	case S:
		return "S"
	case T:
		return "T"
	case U:
		return "U"
	case V:
		return "V"
	case W:
		return "W"
	case X:
		return "X"
	case Y:
		return "Y"
	case Z:
		return "Z"
	case COMMA:
		return ","
	case PERIOD:
		return "."
	case ALT_LEFT:
		return "L-Alt"
	case ALT_RIGHT:
		return "R-Alt"
	case SHIFT_LEFT:
		return "L-Shift"
	case SHIFT_RIGHT:
		return "R-Shift"
	case TAB:
		return "Tab"
	case SPACE:
		return "Space"
	case SYM:
		return "SYM"
	case EXPLORER:
		return "Explorer"
	case ENVELOPE:
		return "Envelope"
	case ENTER:
		return "Enter"
	case DEL:
		return "Delete" // also BACKSPACE
	case GRAVE:
		return "`"
	case MINUS:
		return "-"
	case EQUALS:
		return "="
	case LEFT_BRACKET:
		return "["
	case RIGHT_BRACKET:
		return "]"
	case BACKSLASH:
		return "\\"
	case SEMICOLON:
		return ";"
	case APOSTROPHE:
		return "'"
	case SLASH:
		return "/"
	case AT:
		return "@"
	case NUM:
		return "Num"
	case HEADSETHOOK:
		return "Headset Hook"
	case FOCUS:
		return "Focus"
	case PLUS:
		return "Plus"
	case MENU:
		return "Menu"
	case NOTIFICATION:
		return "Notification"
	case SEARCH:
		return "Search"
	case MEDIA_PLAY_PAUSE:
		return "Play/Pause"
	case MEDIA_STOP:
		return "Stop Media"
	case MEDIA_NEXT:
		return "Next Media"
	case MEDIA_PREVIOUS:
		return "Prev Media"
	case MEDIA_REWIND:
		return "Rewind"
	case MEDIA_FAST_FORWARD:
		return "Fast Forward"
	case MUTE:
		return "Mute"
	case PAGE_UP:
		return "Page Up"
	case PAGE_DOWN:
		return "Page Down"
	case PICTSYMBOLS:
		return "PICTSYMBOLS"
	case SWITCH_CHARSET:
		return "SWITCH_CHARSET"
	case BUTTON_A:
		return "A Button"
	case BUTTON_B:
		return "B Button"
	case BUTTON_C:
		return "C Button"
	case BUTTON_X:
		return "X Button"
	case BUTTON_Y:
		return "Y Button"
	case BUTTON_Z:
		return "Z Button"
	case BUTTON_L1:
		return "L1 Button"
	case BUTTON_R1:
		return "R1 Button"
	case BUTTON_L2:
		return "L2 Button"
	case BUTTON_R2:
		return "R2 Button"
	case BUTTON_THUMBL:
		return "Left Thumb"
	case BUTTON_THUMBR:
		return "Right Thumb"
	case BUTTON_START:
		return "Start"
	case BUTTON_SELECT:
		return "Select"
	case BUTTON_MODE:
		return "Button Mode"
	case FORWARD_DEL:
		return "Forward Delete"
	case CONTROL_LEFT:
		return "L-Ctrl"
	case CONTROL_RIGHT:
		return "R-Ctrl"
	case ESCAPE:
		return "Escape"
	case END:
		return "End"
	case INSERT:
		return "Insert"
	case NUMPAD_0:
		return "Numpad 0"
	case NUMPAD_1:
		return "Numpad 1"
	case NUMPAD_2:
		return "Numpad 2"
	case NUMPAD_3:
		return "Numpad 3"
	case NUMPAD_4:
		return "Numpad 4"
	case NUMPAD_5:
		return "Numpad 5"
	case NUMPAD_6:
		return "Numpad 6"
	case NUMPAD_7:
		return "Numpad 7"
	case NUMPAD_8:
		return "Numpad 8"
	case NUMPAD_9:
		return "Numpad 9"
	case COLON:
		return ":"
	case F1:
		return "F1"
	case F2:
		return "F2"
	case F3:
		return "F3"
	case F4:
		return "F4"
	case F5:
		return "F5"
	case F6:
		return "F6"
	case F7:
		return "F7"
	case F8:
		return "F8"
	case F9:
		return "F9"
	case F10:
		return "F10"
	case F11:
		return "F11"
	case F12:
		return "F12"
		// BUTTON_CIRCLE unhandled, as it conflicts with the more likely to be pressed F12
	default:
		// key name not found
		return ""
	}
}
