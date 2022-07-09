package keyboard

import (
	"bytes"
	"os"
	"time"
)

//File represents the keyboard device file in user space /dev/hidg<#>
type File struct {
	os.File
	StrokeDelay time.Duration
}

/* Keyboard HID Report Descriptor
Modifiers Keys Bit Assignment

  Bit 7     Bit 6   Bit 5    Bit 4    Bit 3	   Bit 2	Bit 1     Bit 0
┏━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━┯━━━━━━━━┓
┃ Right  │ Right  │ Right  │ Right  │ Left   │ Left   │ Left   │ Left   ┃ 1: ON
┃ GUI    │ ALT    │ SHIFT  │ CTRL   │ GUI    │ ALT    │ SHIFT  │ CTRL   ┃ 2: OFF
┗━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━┷━━━━━━━━┛

Byte
	┏━━━━━━━━━━━━━━━┓
0	┃ Modifier Keys ┃
	┠───────────────┨
1	┃ Reserved      ┃
	┠───────────────┨
2	┃ Keycode 1     ┃
	┠───────────────┨
3	┃ Keycode 2     ┃
	┠───────────────┨
4	┃ Keycode 3     ┃
	┠───────────────┨
5	┃ Keycode 4     ┃
	┠───────────────┨
6	┃ Keycode 5     ┃
	┠───────────────┨
7	┃ Keycode 6     ┃
	┗━━━━━━━━━━━━━━━┛
*/
// Report represents a keyboard report sent from the keyboard to the host.
//   Send null '\0\0\0\0\0\0\0\0' report to signify key release.
//   Notes:
//   In playing around with sending keystrokes to Linux -
//     It appears sending a new report clears the old pressed keys? I only needed to send a null report as the last key stroke and things seemed to work.
//     Reports seem to allow from 1 to 6 keycodes. I was able to send a report of size 3 (bytes), one scan code, and it still worked.
const ReportSz int = 8

type Report [ReportSz]byte

func (f *File) WriteString(s string) (n int, err error) {
	var buf bytes.Buffer = bytes.Buffer{}
	for _, c := range s {
		modifier, keycode := ASCII_to_keycode(byte(c))
		r := Report{modifier, 0, keycode, 0, 0, 0, 0, 0}
		buf.Write(r[:])
		r = Report{0, 0, 0, 0, 0, 0, 0, 0}
		buf.Write(r[:])
	}
	return f.File.Write(buf.Bytes())
}

func (f *File) WriteStringDelayed(s string) (n int, err error) {
	var totalBytes int

	for _, c := range s {
		modifier, keycode := ASCII_to_keycode(byte(c))
		r := Report{modifier, 0, keycode, 0, 0, 0, 0, 0}
		if n, err = f.File.Write(r[:]); err != nil {
			totalBytes += n
			return totalBytes, err
		}
		totalBytes += n
		time.Sleep(f.StrokeDelay)
		r = Report{0, 0, 0, 0, 0, 0, 0, 0}
		if n, err = f.File.Write(r[:]); err != nil {
			totalBytes += n
			return totalBytes, err
		}
		totalBytes += n
		time.Sleep(f.StrokeDelay)
	}
	return totalBytes, nil
}

// void pressKey(uint8_t modifiers, uint8_t keycode1, uint8_t keycode2, uint8_t keycode3, uint8_t keycode4, uint8_t keycode5, uint8_t keycode6);
// void pressKey(uint8_t modifiers, uint8_t keycode1, uint8_t keycode2, uint8_t keycode3, uint8_t keycode4, uint8_t keycode5);
// void pressKey(uint8_t modifiers, uint8_t keycode1, uint8_t keycode2, uint8_t keycode3, uint8_t keycode4);
// void pressKey(uint8_t modifiers, uint8_t keycode1, uint8_t keycode2, uint8_t keycode3);
// void pressKey(uint8_t modifiers, uint8_t keycode1, uint8_t keycode2);
// void pressKey(uint8_t modifiers, uint8_t keycode1);
// // presses a list of keys
// void pressKeys(uint8_t modifiers, uint8_t* keycodes, uint8_t sz);

func ASCII_to_keycode(ascii byte) (modifier, keycode byte) {
	keycode, modifier = KEYCODE_NIL, MODIFIER_NOT_SET

	// see scancode.doc appendix C

	if ascii >= 'A' && ascii <= 'Z' {
		keycode = 4 + ascii - 'A'           // set letter
		modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
	} else if ascii >= 'a' && ascii <= 'z' {
		keycode = 4 + ascii - 'a'            // set letter
		modifier &= ^MODIFIER_KEY_LEFT_SHIFT // no shift
	} else if ascii >= '0' && ascii <= '9' {
		modifier = MODIFIER_NOT_SET
		if ascii == '0' {
			keycode = KEYCODE_0
		} else {
			keycode = 30 + ascii - '1'
		}
	} else {
		switch ascii { // convert ascii to keycode according to documentation
		case '!':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			keycode = KEYCODE_Z + 1
		case '@':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			keycode = KEYCODE_Z + 2
		case '#':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			keycode = KEYCODE_Z + 3
		case '$':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			keycode = KEYCODE_Z + 4
		case '%':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			keycode = KEYCODE_Z + 5
		case '^':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			keycode = KEYCODE_Z + 6
		case '&':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			keycode = KEYCODE_Z + 7
		case '*':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			keycode = KEYCODE_Z + 8
		case '(':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			keycode = KEYCODE_Z + 9
		case ')':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			keycode = KEYCODE_0
		case '~':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case '`':
			keycode = KEYCODE_BACK_TICK
		case '_':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case '-':
			keycode = KEYCODE_MINUS
		case '+':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case '=':
			keycode = KEYCODE_EQUAL
		case '{':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case '[':
			keycode = KEYCODE_SQBRAK_LEFT
		case '}':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case ']':
			keycode = KEYCODE_SQBRAK_RIGHT
		case '|':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case '\\':
			keycode = KEYCODE_BACKSLASH
		case ':':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case ';':
			keycode = KEYCODE_SEMICOLON
		case '"':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case '\'':
			keycode = KEYCODE_SINGLE_QUOTE
		case '<':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case ',':
			keycode |= KEYCODE_COMMA
		case '>':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case '.':
			keycode = KEYCODE_PERIOD
		case '?':
			modifier |= MODIFIER_KEY_LEFT_SHIFT // hold shift
			fallthrough
		case '/':
			keycode = KEYCODE_SLASH
		case ' ':
			keycode = KEYCODE_SPACE
		case '\t':
			keycode = KEYCODE_TAB
		case '\n':
			keycode = KEYCODE_ENTER
		}
	}

	return
}

const (
	MODIFIER_KEY_RIGHT_GUI   byte = 0b1000_0000
	MODIFIER_KEY_RIGHT_ALT   byte = 0b0100_0000
	MODIFIER_KEY_RIGHT_SHIFT byte = 0b0010_0000
	MODIFIER_KEY_RIGHT_CTRL  byte = 0b0001_0000
	MODIFIER_KEY_LEFT_GUI    byte = 0b0000_1000
	MODIFIER_KEY_LEFT_ALT    byte = 0b0000_0100
	MODIFIER_KEY_LEFT_SHIFT  byte = 0b0000_0010
	MODIFIER_KEY_LEFT_CTRL   byte = 0b0000_0001
	MODIFIER_NOT_SET         byte = 0b0000_0000
)

const (
	// some convenience definitions for modifier keys
	KEYCODE_NIL               byte = 0x00
	KEYCODE_MOD_LEFT_CONTROL  byte = 0x01
	KEYCODE_MOD_LEFT_SHIFT    byte = 0x02
	KEYCODE_MOD_LEFT_ALT      byte = 0x04
	KEYCODE_MOD_LEFT_GUI      byte = 0x08
	KEYCODE_MOD_RIGHT_CONTROL byte = 0x10
	KEYCODE_MOD_RIGHT_SHIFT   byte = 0x20
	KEYCODE_MOD_RIGHT_ALT     byte = 0x40
	KEYCODE_MOD_RIGHT_GUI     byte = 0x80

	// some more keycodes
	KEYCODE_LEFT_CONTROL byte = 0xE0
	KEYCODE_LEFT_SHIFT   byte = 0xE1
	KEYCODE_LEFT_ALT     byte = 0xE2
	KEYCODE_LEFT_GUI     byte = 0xE3
	KEYCODE_RIGHT_CONTRO byte = 0xE4
	KEYCODE_RIGHT_SHIFT  byte = 0xE5
	KEYCODE_RIGHT_ALT    byte = 0xE6
	KEYCODE_RIGHT_GUI    byte = 0xE7
	KEYCODE_1            byte = 0x1E
	KEYCODE_2            byte = 0x1F
	KEYCODE_3            byte = 0x20
	KEYCODE_4            byte = 0x21
	KEYCODE_5            byte = 0x22
	KEYCODE_6            byte = 0x23
	KEYCODE_7            byte = 0x24
	KEYCODE_8            byte = 0x25
	KEYCODE_9            byte = 0x26
	KEYCODE_0            byte = 0x27
	KEYCODE_A            byte = 0x04
	KEYCODE_B            byte = 0x05
	KEYCODE_C            byte = 0x06
	KEYCODE_D            byte = 0x07
	KEYCODE_E            byte = 0x08
	KEYCODE_F            byte = 0x09
	KEYCODE_G            byte = 0x0A
	KEYCODE_H            byte = 0x0B
	KEYCODE_I            byte = 0x0C
	KEYCODE_J            byte = 0x0D
	KEYCODE_K            byte = 0x0E
	KEYCODE_L            byte = 0x0F
	KEYCODE_M            byte = 0x10
	KEYCODE_N            byte = 0x11
	KEYCODE_O            byte = 0x12
	KEYCODE_P            byte = 0x13
	KEYCODE_Q            byte = 0x14
	KEYCODE_R            byte = 0x15
	KEYCODE_S            byte = 0x16
	KEYCODE_T            byte = 0x17
	KEYCODE_U            byte = 0x18
	KEYCODE_V            byte = 0x19
	KEYCODE_W            byte = 0x1A
	KEYCODE_X            byte = 0x1B
	KEYCODE_Y            byte = 0x1C
	KEYCODE_Z            byte = 0x1D
	KEYCODE_COMMA        byte = 0x36
	KEYCODE_PERIOD       byte = 0x37
	KEYCODE_MINUS        byte = 0x2D
	KEYCODE_EQUAL        byte = 0x2E
	KEYCODE_BACKSLASH    byte = 0x31
	KEYCODE_SQBRAK_LEFT  byte = 0x2F
	KEYCODE_SQBRAK_RIGHT byte = 0x30
	KEYCODE_COLON        byte = 0x33
	KEYCODE_SEMICOLON    byte = 0x33
	KEYCODE_SINGLE_QUOTE byte = 0x34
	KEYCODE_DOUBLE_QUOTE byte = 0x34
	KEYCODE_BACK_TICK    byte = 0x35
	KEYCODE_TILDA        byte = 0x35
	KEYCODE_SLASH        byte = 0x38
	KEYCODE_F1           byte = 0x3A
	KEYCODE_F2           byte = 0x3B
	KEYCODE_F3           byte = 0x3C
	KEYCODE_F4           byte = 0x3D
	KEYCODE_F5           byte = 0x3E
	KEYCODE_F6           byte = 0x3F
	KEYCODE_F7           byte = 0x40
	KEYCODE_F8           byte = 0x41
	KEYCODE_F9           byte = 0x42
	KEYCODE_F10          byte = 0x43
	KEYCODE_F11          byte = 0x44
	KEYCODE_F12          byte = 0x45
	KEYCODE_APP          byte = 0x65
	KEYCODE_ENTER        byte = 0x28
	KEYCODE_BACKSPACE    byte = 0x2A
	KEYCODE_ESC          byte = 0x29
	KEYCODE_TAB          byte = 0x2B
	KEYCODE_SPACE        byte = 0x2C
	KEYCODE_INSERT       byte = 0x49
	KEYCODE_HOME         byte = 0x4A
	KEYCODE_PAGE_UP      byte = 0x4B
	KEYCODE_DELETE       byte = 0x4C
	KEYCODE_END          byte = 0x4D
	KEYCODE_PAGE_DOWN    byte = 0x4E
	KEYCODE_PRINTSCREEN  byte = 0x46
	KEYCODE_ARROW_RIGHT  byte = 0x4F
	KEYCODE_ARROW_LEFT   byte = 0x50
	KEYCODE_ARROW_DOWN   byte = 0x51
	KEYCODE_ARROW_UP     byte = 0x52
)
