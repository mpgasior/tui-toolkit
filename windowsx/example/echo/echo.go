//go:build windows

package main

import (
	"fmt"
	"os"

	"github.com/nimelo/tui-go/utf16x"
	"github.com/nimelo/tui-go/windowsx"
	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

func main() {
	fd := os.Stdin.Fd()

	state, _ := term.MakeRaw(int(fd))
	defer term.Restore(int(fd), state)

	d := utf16x.Decoder{}
	buffer := make([]windowsx.INPUT_RECORD, 1)
	for {
		_, err := windowsx.ReadConsoleInput(windows.Handle(fd), buffer)
		if err != nil {
			panic(err)
		}

		record := buffer[0]
		if keyEvent, ok := record.KeyEvent(); ok {
			if keyEvent.KeyDown == 0 {
				continue
			}

			if keyEvent.VirtualKeyCode != 0 && keyEvent.UnicodeChar == 0 {
				continue
			}

			r, ok := d.Decode(keyEvent.UnicodeChar)
			if !ok {
				continue
			}

			if r == '\x03' {
				break
			}

			fmt.Printf("%c", r)
		}
	}
}
