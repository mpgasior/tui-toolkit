package main

import (
	"fmt"
	"os"

	"github.com/nimelo/tui-go/windowsx"
	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

func main() {
	fd := os.Stdin.Fd()

	state, err := term.MakeRaw(int(fd))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(fd), state)

	buffer := make([]windowsx.INPUT_RECORD, 1)
	for {
		_, err := windowsx.ReadConsoleInput(windows.Handle(fd), buffer)
		if err != nil {
			panic(err)
		}

		record := buffer[0]
		fmt.Printf("%s\n", stringify(record))

		if record.EventType == windowsx.KEY_EVENT {
			keyEvent := record.KeyEvent()

			if keyEvent.Char[0] == 'q' {
				break
			}
		}
	}
}

func stringify(r windowsx.INPUT_RECORD) string {
	switch r.EventType {
	case windowsx.KEY_EVENT:
		k := r.KeyEvent()
		state := "UP"
		if k.KeyDown != 0 {
			state = "DOWN"
		}
		return fmt.Sprintf("KEY_EVENT: [%s] VK: 0x%02X, Char: %c", state, k.VirtualKeyCode, k.Char[0])

	case windowsx.MOUSE_EVENT:
		m := r.MouseEvent()
		return fmt.Sprintf("MOUSE_EVENT: Pos(%d, %d), Buttons: 0x%X", m.MousePosition.X, m.MousePosition.Y, m.ButtonState)

	case windowsx.WINDOW_BUFFER_SIZE_EVENT:
		s := r.WindowsBufferSizeEvent()
		return fmt.Sprintf("RESIZE: %dx%d", s.Size.X, s.Size.Y)

	default:
		return fmt.Sprintf("EVENT_TYPE: %d", r.EventType)
	}
}
