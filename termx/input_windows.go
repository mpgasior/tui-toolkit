//go:build windows

package termx

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/nimelo/tui-go/windowsx"
	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

type terminalInput struct {
	f         *os.File
	stopEvent windows.Handle
}

func NewTerminalInput(f *os.File) (TerminalInput, error) {
	event, err := windows.CreateEvent(nil, 1, 0, nil)

	if err != nil {
		return nil, err
	}

	r := &terminalInput{
		f:         f,
		stopEvent: event,
	}

	return r, nil
}

func (ti *terminalInput) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	handle := windows.Handle(ti.f.Fd())
	windows.ResetEvent(ti.stopEvent)

	stop := context.AfterFunc(ctx, func() {
		windows.SetEvent(ti.stopEvent)
	})
	defer stop()

	handles := []windows.Handle{handle, ti.stopEvent}

	event, err := windows.WaitForMultipleObjects(handles, false, windows.INFINITE)
	if err != nil {
		return 0, err
	}
	switch event {
	case windows.WAIT_OBJECT_0:
		buffer := make([]windowsx.INPUT_RECORD, 1)

		_, err := windowsx.ReadConsoleInput(handle, buffer)
		if err != nil {
			return 0, nil
		}

		r := buffer[0]
		if r.EventType == windowsx.KEY_EVENT {
			keyEvent := r.KeyEvent()

			if keyEvent.KeyDown == 1 {
				p[0] = keyEvent.Char[0]

				return 1, nil
			}
		}

		return 0, nil

	case windows.WAIT_OBJECT_0 + 1:
		return 0, ctx.Err()
	case windows.WAIT_FAILED:
		return 0, windows.GetLastError()
	default:
		return 0, errors.New("unexpected wait result")
	}
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

func (ti *terminalInput) Close() error {
	err := errors.Join(
		windows.SetEvent(ti.stopEvent),
		windows.CloseHandle(ti.stopEvent))

	return err
}
