//go:build windows

package main

import (
	"fmt"
	"os"
	"unicode"

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
		printRecord(record)

		if keyEvent, ok := record.KeyEvent(); ok {
			if keyEvent.UnicodeChar == 'q' {
				break
			}
		}
	}
}

func printRecord(r windowsx.INPUT_RECORD) {
	fmt.Printf("%s: ", r.EventType.String())

	if e, ok := r.KeyEvent(); ok {
		printKey(e)
	}

	if e, ok := r.FocusEvent(); ok {
		printFocus(e)
	}

	if e, ok := r.MouseEvent(); ok {
		printMouse(e)
	}

	if e, ok := r.WindowsBufferSizeEvent(); ok {
		printWindow(e)
	}

	fmt.Print("\r\n")
}

func printKey(e *windowsx.KEY_EVENT_RECORD) {
	state := "UP"
	if e.KeyDown != 0 {
		state = "DOWN"
	}

	source := "NATIVE"
	if e.VirtualKeyCode == 0 {
		source = "SYNTH"
	}

	fmt.Printf("[%s] ", state)
	fmt.Printf("[%s] ", source)
	fmt.Printf("%s ", e.VirtualKeyCode.String())
	fmt.Printf("%s ", e.ControlKeyState.String())
	if unicode.IsPrint(rune(e.UnicodeChar)) {
		fmt.Printf("UnicodeChar: '%c' (0x%02X) ", e.UnicodeChar, e.UnicodeChar)
	} else {
		fmt.Printf("UnicodeChar: (0x%02X) ", e.UnicodeChar)
	}
	fmt.Printf("Count: %d ", e.RepeatCount)
}

func printWindow(s *windowsx.WINDOW_BUFFER_SIZE_RECORD) {
	fmt.Printf("%dx%d", s.Size.X, s.Size.Y)
}

func printMouse(m *windowsx.MOUSE_EVENT_RECORD) {
	fmt.Printf("Pos(%d, %d), Buttons: %s", m.MousePosition.X, m.MousePosition.Y, m.ButtonState.String())
}

func printFocus(e *windowsx.FOCUS_EVENT_RECORD) {
	focus := "ENTER"
	if e.SetFocus == 0 {
		focus = "EXIT"
	}

	fmt.Printf("[%s]", focus)
}
