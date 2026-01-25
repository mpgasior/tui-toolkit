//go:build windows

package termx

import (
	"bytes"
	"context"
	"errors"
	"os"
	"sync"
	"unicode/utf8"

	"github.com/mpgasior/tui-go/utf16x"
	"github.com/mpgasior/tui-go/windowsx"
	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

type terminalInput struct {
	f         *os.File
	stopEvent windows.Handle
	decoder   utf16x.Decoder
	buffer    bytes.Buffer
}

func NewTerminalInput(f *os.File) (*terminalInput, error) {
	event, err := windows.CreateEvent(nil, 0, 0, nil)

	if err != nil {
		return nil, err
	}

	r := &terminalInput{
		f:         f,
		stopEvent: event,
	}

	return r, nil
}

func (ti *terminalInput) MakeRaw() (func() error, error) {
	fd := int(ti.f.Fd())
	state, err := term.MakeRaw(fd)

	if err != nil {
		return nil, err
	}

	var once sync.Once
	var restoreErr error

	restore := func() error {
		once.Do(func() {
			restoreErr = term.Restore(fd, state)
		})
		return restoreErr
	}

	return restore, nil
}

func (ti *terminalInput) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	if ti.buffer.Len() > 0 {
		return ti.buffer.Read(p)
	}

	console := windows.Handle(ti.f.Fd())
	buffer := make([]windowsx.INPUT_RECORD, 1024)

	for {
		n, err := windowsx.PeekConsoleInput(console, buffer)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			if err := ti.waitEvent(ctx); err != nil {
				return 0, err
			}
			continue
		}

		n, err = windowsx.ReadConsoleInput(console, buffer)
		if err != nil {
			return 0, err
		}

		for i := range n {
			record := buffer[i]

			keyEvent, ok := record.KeyEvent()
			if !ok || keyEvent.KeyDown == 0 {
				continue
			}

			if keyEvent.VirtualKeyCode != 0 && keyEvent.UnicodeChar == 0 {
				continue
			}

			r, ok := ti.decoder.Decode(keyEvent.UnicodeChar)
			if !ok {
				continue
			}

			var runeBytes [utf8.UTFMax]byte
			runeLength := utf8.EncodeRune(runeBytes[:], r)
			ti.buffer.Write(runeBytes[:runeLength])
		}

		if ti.buffer.Len() > 0 {
			return ti.buffer.Read(p)
		}
	}
}

func (ti *terminalInput) waitEvent(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	console := windows.Handle(ti.f.Fd())

	stop := context.AfterFunc(ctx, func() {
		windows.SetEvent(ti.stopEvent)
	})
	defer stop()

	handles := []windows.Handle{console, ti.stopEvent}
	event, err := windows.WaitForMultipleObjects(handles, false, windows.INFINITE)
	if err != nil {
		return err
	}

	switch event {
	case windows.WAIT_OBJECT_0:
		return nil
	case windows.WAIT_OBJECT_0 + 1:
		return ctx.Err()
	case windows.WAIT_FAILED:
		return windows.GetLastError()
	default:
		return errors.New("unexpected wait result")
	}
}

func (ti *terminalInput) Close() error {
	err := errors.Join(
		windows.SetEvent(ti.stopEvent),
		windows.CloseHandle(ti.stopEvent))

	return err
}
