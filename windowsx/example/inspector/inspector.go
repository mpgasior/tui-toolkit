//go:build windows

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mpgasior/tui-go/windowsx"
	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

func main() {
	fd := os.Stdin.Fd()

	state, _ := term.MakeRaw(int(fd))
	defer term.Restore(int(fd), state)

	buffer := make([]windowsx.INPUT_RECORD, 1)
	for {
		_, err := windowsx.ReadConsoleInput(windows.Handle(fd), buffer)
		if err != nil {
			panic(err)
		}

		record := buffer[0]
		fmt.Printf("%s\n", stringify(record))

		if keyEvent, ok := record.KeyEvent(); ok {
			if keyEvent.UnicodeChar == 'q' {
				break
			}
		}
	}
}

func stringify(r windowsx.INPUT_RECORD) string {
	if e, ok := r.KeyEvent(); ok {
		return stringifyKey(e)
	}

	if e, ok := r.FocusEvent(); ok {
		stringifyFocus(e)
	}

	if e, ok := r.MouseEvent(); ok {
		return stringifyMouse(e)
	}

	if e, ok := r.WindowsBufferSizeEvent(); ok {
		return stringifyWindow(e)
	}

	return fmt.Sprintf("EVENT_TYPE: %d", r.EventType)
}

func stringifyKey(e *windowsx.KEY_EVENT_RECORD) string {
	state := "UP"
	if e.KeyDown != 0 {
		state = "DOWN"
	}

	source := "NATIVE"
	if e.VirtualKeyCode == 0 {
		source = "SYNTH"
	}

	vk, ok := windowsx.VirtualKeyMap[e.VirtualKeyCode]
	if !ok {
		vk = "-"
	}

	var sb strings.Builder

	sb.WriteString("KEY_EVENT: ")
	sb.WriteString(fmt.Sprintf("[%s] ", state))
	sb.WriteString(fmt.Sprintf("[%s] ", source))
	sb.WriteString(fmt.Sprintf("VK: %s (0x%02X) ", vk, e.VirtualKeyCode))
	sb.WriteString(fmt.Sprintf("ControlKey: %s ", stringifyControlKey(e.ControlKeyState)))
	if e.UnicodeChar >= 32 && e.UnicodeChar <= 126 {
		sb.WriteString(fmt.Sprintf("UnicodeChar: '%c' (0x%02X) ", e.UnicodeChar, e.UnicodeChar))
	} else {
		sb.WriteString(fmt.Sprintf("UnicodeChar: [CTRL](0x%02X) ", e.UnicodeChar))
	}
	sb.WriteString(fmt.Sprintf("Count: %d ", e.RepeatCount))

	return sb.String()
}

func stringifyControlKey(state uint32) string {
	var sb strings.Builder

	first := true
	for k, v := range windowsx.ControlKeyStateMap {
		if state&k != 0 {
			if !first {
				sb.WriteString(" | ")
			}
			first = false

			sb.WriteString(v)
		}
	}

	if sb.Len() == 0 {
		return "-"
	}

	return sb.String()
}

func stringifyWindow(s *windowsx.WINDOW_BUFFER_SIZE_RECORD) string {
	return fmt.Sprintf("RESIZE: %dx%d", s.Size.X, s.Size.Y)
}

func stringifyMouse(m *windowsx.MOUSE_EVENT_RECORD) string {
	return fmt.Sprintf("MOUSE_EVENT: Pos(%d, %d), Buttons: 0x%X", m.MousePosition.X, m.MousePosition.Y, m.ButtonState)
}

func stringifyFocus(e *windowsx.FOCUS_EVENT_RECORD) string {
	focus := "ENTER"
	if e.SetFocus == 0 {
		focus = "EXIT"
	}

	return fmt.Sprintf("FOCUS_EVENT: [%s]", focus)
}
