//go:build windows

package termx

import (
	"context"
	"errors"
	"os"

	"golang.org/x/sys/windows"
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
		return ti.f.Read(p)
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

	handle := windows.Handle(ti.f.Fd())
	if err := windows.FlushConsoleInputBuffer(handle); err != nil {
		return nil, err
	}

	return makeRaw(fd)
}

func (ti *terminalInput) Close() error {
	err := errors.Join(
		windows.SetEvent(ti.stopEvent),
		windows.CloseHandle(ti.stopEvent))

	return err
}
