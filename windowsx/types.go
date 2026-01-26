//go:build windows

package windowsx

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=EventType,VirtualKeyCode -output=types_string_windows.go
//go:generate go run github.com/mpgasior/tui-go/tools/bitstringer -type=ControlKeyState,MouseButtonState,MouseEventFlag -output=types_bitstring_windows.go

type EventType uint16

const (
	KEY_EVENT                EventType = 0x0001
	MOUSE_EVENT              EventType = 0x0002
	WINDOW_BUFFER_SIZE_EVENT EventType = 0x0004
	MENU_EVENT               EventType = 0x0008
	FOCUS_EVENT              EventType = 0x0010
)

// See: https://learn.microsoft.com/en-us/windows/console/input-record-str
type INPUT_RECORD struct {
	EventType EventType
	_         uint16
	Event     [16]byte
}

// See: https://learn.microsoft.com/en-us/windows/console/focus-event-record-str
type FOCUS_EVENT_RECORD struct {
	SetFocus uint32
}

type ControlKeyState uint32

const (
	CAPSLOCK_ON        ControlKeyState = 0x0080
	ENHANCED_KEY       ControlKeyState = 0x0100
	LEFT_ALT_PRESSED   ControlKeyState = 0x0002
	LEFT_CTRL_PRESSED  ControlKeyState = 0x0008
	NUMLOCK_ON         ControlKeyState = 0x0020
	RIGHT_ALT_PRESSED  ControlKeyState = 0x0001
	RIGHT_CTRL_PRESSED ControlKeyState = 0x0004
	SCROLLLOCK_ON      ControlKeyState = 0x0040
	SHIFT_PRESSED      ControlKeyState = 0x0010
)

type VirtualKeyCode uint16

// See https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
const (
	VK_LBUTTON                         VirtualKeyCode = 0x01
	VK_RBUTTON                         VirtualKeyCode = 0x02
	VK_CANCEL                          VirtualKeyCode = 0x03
	VK_MBUTTON                         VirtualKeyCode = 0x04
	VK_XBUTTON1                        VirtualKeyCode = 0x05
	VK_XBUTTON2                        VirtualKeyCode = 0x06
	VK_BACK                            VirtualKeyCode = 0x08
	VK_TAB                             VirtualKeyCode = 0x09
	VK_CLEAR                           VirtualKeyCode = 0x0C
	VK_RETURN                          VirtualKeyCode = 0x0D
	VK_SHIFT                           VirtualKeyCode = 0x10
	VK_CONTROL                         VirtualKeyCode = 0x11
	VK_MENU                            VirtualKeyCode = 0x12
	VK_PAUSE                           VirtualKeyCode = 0x13
	VK_CAPITAL                         VirtualKeyCode = 0x14
	VK_KANA                            VirtualKeyCode = 0x15
	VK_HANGUL                          VirtualKeyCode = 0x15
	VK_IME_ON                          VirtualKeyCode = 0x16
	VK_JUNJA                           VirtualKeyCode = 0x17
	VK_FINAL                           VirtualKeyCode = 0x18
	VK_HANJA                           VirtualKeyCode = 0x19
	VK_KANJI                           VirtualKeyCode = 0x19
	VK_IME_OFF                         VirtualKeyCode = 0x1A
	VK_ESCAPE                          VirtualKeyCode = 0x1B
	VK_CONVERT                         VirtualKeyCode = 0x1C
	VK_NONCONVERT                      VirtualKeyCode = 0x1D
	VK_ACCEPT                          VirtualKeyCode = 0x1E
	VK_MODECHANGE                      VirtualKeyCode = 0x1F
	VK_SPACE                           VirtualKeyCode = 0x20
	VK_PRIOR                           VirtualKeyCode = 0x21
	VK_NEXT                            VirtualKeyCode = 0x22
	VK_END                             VirtualKeyCode = 0x23
	VK_HOME                            VirtualKeyCode = 0x24
	VK_LEFT                            VirtualKeyCode = 0x25
	VK_UP                              VirtualKeyCode = 0x26
	VK_RIGHT                           VirtualKeyCode = 0x27
	VK_DOWN                            VirtualKeyCode = 0x28
	VK_SELECT                          VirtualKeyCode = 0x29
	VK_PRINT                           VirtualKeyCode = 0x2A
	VK_EXECUTE                         VirtualKeyCode = 0x2B
	VK_SNAPSHOT                        VirtualKeyCode = 0x2C
	VK_INSERT                          VirtualKeyCode = 0x2D
	VK_DELETE                          VirtualKeyCode = 0x2E
	VK_HELP                            VirtualKeyCode = 0x2F
	VK_0                               VirtualKeyCode = 0x30
	VK_1                               VirtualKeyCode = 0x31
	VK_2                               VirtualKeyCode = 0x32
	VK_3                               VirtualKeyCode = 0x33
	VK_4                               VirtualKeyCode = 0x34
	VK_5                               VirtualKeyCode = 0x35
	VK_6                               VirtualKeyCode = 0x36
	VK_7                               VirtualKeyCode = 0x37
	VK_8                               VirtualKeyCode = 0x38
	VK_9                               VirtualKeyCode = 0x39
	VK_A                               VirtualKeyCode = 0x41
	VK_B                               VirtualKeyCode = 0x42
	VK_C                               VirtualKeyCode = 0x43
	VK_D                               VirtualKeyCode = 0x44
	VK_E                               VirtualKeyCode = 0x45
	VK_F                               VirtualKeyCode = 0x46
	VK_G                               VirtualKeyCode = 0x47
	VK_H                               VirtualKeyCode = 0x48
	VK_I                               VirtualKeyCode = 0x49
	VK_J                               VirtualKeyCode = 0x4A
	VK_K                               VirtualKeyCode = 0x4B
	VK_L                               VirtualKeyCode = 0x4C
	VK_M                               VirtualKeyCode = 0x4D
	VK_N                               VirtualKeyCode = 0x4E
	VK_O                               VirtualKeyCode = 0x4F
	VK_P                               VirtualKeyCode = 0x50
	VK_Q                               VirtualKeyCode = 0x51
	VK_R                               VirtualKeyCode = 0x52
	VK_S                               VirtualKeyCode = 0x53
	VK_T                               VirtualKeyCode = 0x54
	VK_U                               VirtualKeyCode = 0x55
	VK_V                               VirtualKeyCode = 0x56
	VK_W                               VirtualKeyCode = 0x57
	VK_X                               VirtualKeyCode = 0x58
	VK_Y                               VirtualKeyCode = 0x59
	VK_Z                               VirtualKeyCode = 0x5A
	VK_LWIN                            VirtualKeyCode = 0x5B
	VK_RWIN                            VirtualKeyCode = 0x5C
	VK_APPS                            VirtualKeyCode = 0x5D
	VK_SLEEP                           VirtualKeyCode = 0x5F
	VK_NUMPAD0                         VirtualKeyCode = 0x60
	VK_NUMPAD1                         VirtualKeyCode = 0x61
	VK_NUMPAD2                         VirtualKeyCode = 0x62
	VK_NUMPAD3                         VirtualKeyCode = 0x63
	VK_NUMPAD4                         VirtualKeyCode = 0x64
	VK_NUMPAD5                         VirtualKeyCode = 0x65
	VK_NUMPAD6                         VirtualKeyCode = 0x66
	VK_NUMPAD7                         VirtualKeyCode = 0x67
	VK_NUMPAD8                         VirtualKeyCode = 0x68
	VK_NUMPAD9                         VirtualKeyCode = 0x69
	VK_MULTIPLY                        VirtualKeyCode = 0x6A
	VK_ADD                             VirtualKeyCode = 0x6B
	VK_SEPARATOR                       VirtualKeyCode = 0x6C
	VK_SUBTRACT                        VirtualKeyCode = 0x6D
	VK_DECIMAL                         VirtualKeyCode = 0x6E
	VK_DIVIDE                          VirtualKeyCode = 0x6F
	VK_F1                              VirtualKeyCode = 0x70
	VK_F2                              VirtualKeyCode = 0x71
	VK_F3                              VirtualKeyCode = 0x72
	VK_F4                              VirtualKeyCode = 0x73
	VK_F5                              VirtualKeyCode = 0x74
	VK_F6                              VirtualKeyCode = 0x75
	VK_F7                              VirtualKeyCode = 0x76
	VK_F8                              VirtualKeyCode = 0x77
	VK_F9                              VirtualKeyCode = 0x78
	VK_F10                             VirtualKeyCode = 0x79
	VK_F11                             VirtualKeyCode = 0x7A
	VK_F12                             VirtualKeyCode = 0x7B
	VK_F13                             VirtualKeyCode = 0x7C
	VK_F14                             VirtualKeyCode = 0x7D
	VK_F15                             VirtualKeyCode = 0x7E
	VK_F16                             VirtualKeyCode = 0x7F
	VK_F17                             VirtualKeyCode = 0x80
	VK_F18                             VirtualKeyCode = 0x81
	VK_F19                             VirtualKeyCode = 0x82
	VK_F20                             VirtualKeyCode = 0x83
	VK_F21                             VirtualKeyCode = 0x84
	VK_F22                             VirtualKeyCode = 0x85
	VK_F23                             VirtualKeyCode = 0x86
	VK_F24                             VirtualKeyCode = 0x87
	VK_NUMLOCK                         VirtualKeyCode = 0x90
	VK_SCROLL                          VirtualKeyCode = 0x91
	VK_LSHIFT                          VirtualKeyCode = 0xA0
	VK_RSHIFT                          VirtualKeyCode = 0xA1
	VK_LCONTROL                        VirtualKeyCode = 0xA2
	VK_RCONTROL                        VirtualKeyCode = 0xA3
	VK_LMENU                           VirtualKeyCode = 0xA4
	VK_RMENU                           VirtualKeyCode = 0xA5
	VK_BROWSER_BACK                    VirtualKeyCode = 0xA6
	VK_BROWSER_FORWARD                 VirtualKeyCode = 0xA7
	VK_BROWSER_REFRESH                 VirtualKeyCode = 0xA8
	VK_BROWSER_STOP                    VirtualKeyCode = 0xA9
	VK_BROWSER_SEARCH                  VirtualKeyCode = 0xAA
	VK_BROWSER_FAVORITES               VirtualKeyCode = 0xAB
	VK_BROWSER_HOME                    VirtualKeyCode = 0xAC
	VK_VOLUME_MUTE                     VirtualKeyCode = 0xAD
	VK_VOLUME_DOWN                     VirtualKeyCode = 0xAE
	VK_VOLUME_UP                       VirtualKeyCode = 0xAF
	VK_MEDIA_NEXT_TRACK                VirtualKeyCode = 0xB0
	VK_MEDIA_PREV_TRACK                VirtualKeyCode = 0xB1
	VK_MEDIA_STOP                      VirtualKeyCode = 0xB2
	VK_MEDIA_PLAY_PAUSE                VirtualKeyCode = 0xB3
	VK_LAUNCH_MAIL                     VirtualKeyCode = 0xB4
	VK_LAUNCH_MEDIA_SELECT             VirtualKeyCode = 0xB5
	VK_LAUNCH_APP1                     VirtualKeyCode = 0xB6
	VK_LAUNCH_APP2                     VirtualKeyCode = 0xB7
	VK_OEM_1                           VirtualKeyCode = 0xBA
	VK_OEM_PLUS                        VirtualKeyCode = 0xBB
	VK_OEM_COMMA                       VirtualKeyCode = 0xBC
	VK_OEM_MINUS                       VirtualKeyCode = 0xBD
	VK_OEM_PERIOD                      VirtualKeyCode = 0xBE
	VK_OEM_2                           VirtualKeyCode = 0xBF
	VK_OEM_3                           VirtualKeyCode = 0xC0
	VK_GAMEPAD_A                       VirtualKeyCode = 0xC3
	VK_GAMEPAD_B                       VirtualKeyCode = 0xC4
	VK_GAMEPAD_X                       VirtualKeyCode = 0xC5
	VK_GAMEPAD_Y                       VirtualKeyCode = 0xC6
	VK_GAMEPAD_RIGHT_SHOULDER          VirtualKeyCode = 0xC7
	VK_GAMEPAD_LEFT_SHOULDER           VirtualKeyCode = 0xC8
	VK_GAMEPAD_LEFT_TRIGGER            VirtualKeyCode = 0xC9
	VK_GAMEPAD_RIGHT_TRIGGER           VirtualKeyCode = 0xCA
	VK_GAMEPAD_DPAD_UP                 VirtualKeyCode = 0xCB
	VK_GAMEPAD_DPAD_DOWN               VirtualKeyCode = 0xCC
	VK_GAMEPAD_DPAD_LEFT               VirtualKeyCode = 0xCD
	VK_GAMEPAD_DPAD_RIGHT              VirtualKeyCode = 0xCE
	VK_GAMEPAD_MENU                    VirtualKeyCode = 0xCF
	VK_GAMEPAD_VIEW                    VirtualKeyCode = 0xD0
	VK_GAMEPAD_LEFT_THUMBSTICK_BUTTON  VirtualKeyCode = 0xD1
	VK_GAMEPAD_RIGHT_THUMBSTICK_BUTTON VirtualKeyCode = 0xD2
	VK_GAMEPAD_LEFT_THUMBSTICK_UP      VirtualKeyCode = 0xD3
	VK_GAMEPAD_LEFT_THUMBSTICK_DOWN    VirtualKeyCode = 0xD4
	VK_GAMEPAD_LEFT_THUMBSTICK_RIGHT   VirtualKeyCode = 0xD5
	VK_GAMEPAD_LEFT_THUMBSTICK_LEFT    VirtualKeyCode = 0xD6
	VK_GAMEPAD_RIGHT_THUMBSTICK_UP     VirtualKeyCode = 0xD7
	VK_GAMEPAD_RIGHT_THUMBSTICK_DOWN   VirtualKeyCode = 0xD8
	VK_GAMEPAD_RIGHT_THUMBSTICK_RIGHT  VirtualKeyCode = 0xD9
	VK_GAMEPAD_RIGHT_THUMBSTICK_LEFT   VirtualKeyCode = 0xDA
	VK_OEM_4                           VirtualKeyCode = 0xDB
	VK_OEM_5                           VirtualKeyCode = 0xDC
	VK_OEM_6                           VirtualKeyCode = 0xDD
	VK_OEM_7                           VirtualKeyCode = 0xDE
	VK_OEM_8                           VirtualKeyCode = 0xDF
	VK_OEM_102                         VirtualKeyCode = 0xE2
	VK_PROCESSKEY                      VirtualKeyCode = 0xE5
	VK_PACKET                          VirtualKeyCode = 0xE7
	VK_ATTN                            VirtualKeyCode = 0xF6
	VK_CRSEL                           VirtualKeyCode = 0xF7
	VK_EXSEL                           VirtualKeyCode = 0xF8
	VK_EREOF                           VirtualKeyCode = 0xF9
	VK_PLAY                            VirtualKeyCode = 0xFA
	VK_ZOOM                            VirtualKeyCode = 0xFB
	VK_NONAME                          VirtualKeyCode = 0xFC
	VK_PA1                             VirtualKeyCode = 0xFD
	VK_OEM_CLEAR                       VirtualKeyCode = 0xFE
)

// See: https://learn.microsoft.com/en-us/windows/console/key-event-record-str
type KEY_EVENT_RECORD struct {
	KeyDown         uint32
	RepeatCount     uint16
	VirtualKeyCode  VirtualKeyCode
	VirtualScanCode uint16
	UnicodeChar     uint16
	ControlKeyState ControlKeyState
}

// See: https://learn.microsoft.com/en-us/windows/console/menu-event-record-str
type MENU_EVENT_RECORD struct {
	CommandId uint32
}

type MouseButtonState uint32

const (
	FROM_LEFT_1ST_BUTTON_PRESSED MouseButtonState = 0x0001
	FROM_LEFT_2ND_BUTTON_PRESSED MouseButtonState = 0x0004
	FROM_LEFT_3RD_BUTTON_PRESSED MouseButtonState = 0x0008
	FROM_LEFT_4TH_BUTTON_PRESSED MouseButtonState = 0x0010
	RIGHTMOST_BUTTON_PRESSED     MouseButtonState = 0x0002
)

type MouseEventFlag uint32

const (
	MOUSE_MOVED    MouseEventFlag = 0x0001
	DOUBLE_CLICK   MouseEventFlag = 0x0002
	MOUSE_WHEELED  MouseEventFlag = 0x0004
	MOUSE_HWHEELED MouseEventFlag = 0x0008
)

// See: https://learn.microsoft.com/en-us/windows/console/mouse-event-record-str
type MOUSE_EVENT_RECORD struct {
	MousePosition   windows.Coord
	ButtonState     MouseButtonState
	ControlKeyState ControlKeyState
	EventFlags      MouseEventFlag
}

// See: https://learn.microsoft.com/en-us/windows/console/window-buffer-size-record-str
type WINDOW_BUFFER_SIZE_RECORD struct {
	Size windows.Coord
}

func (r *INPUT_RECORD) KeyEvent() (*KEY_EVENT_RECORD, bool) {
	if r.EventType == KEY_EVENT {
		return (*KEY_EVENT_RECORD)(unsafe.Pointer(&r.Event[0])), true
	}

	return nil, false
}

func (r *INPUT_RECORD) MouseEvent() (*MOUSE_EVENT_RECORD, bool) {
	if r.EventType == MOUSE_EVENT {
		return (*MOUSE_EVENT_RECORD)(unsafe.Pointer(&r.Event[0])), true
	}

	return nil, false
}

func (r *INPUT_RECORD) WindowsBufferSizeEvent() (*WINDOW_BUFFER_SIZE_RECORD, bool) {
	if r.EventType == WINDOW_BUFFER_SIZE_EVENT {
		return (*WINDOW_BUFFER_SIZE_RECORD)(unsafe.Pointer(&r.Event[0])), true
	}

	return nil, false
}

func (r *INPUT_RECORD) FocusEvent() (*FOCUS_EVENT_RECORD, bool) {
	if r.EventType == FOCUS_EVENT {
		return (*FOCUS_EVENT_RECORD)(unsafe.Pointer(&r.Event[0])), true
	}

	return nil, false
}

func (r *INPUT_RECORD) MenuEvent() (*MENU_EVENT_RECORD, bool) {
	if r.EventType == MENU_EVENT {
		return (*MENU_EVENT_RECORD)(unsafe.Pointer(&r.Event[0])), true
	}

	return nil, false
}
