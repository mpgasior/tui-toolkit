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

type Input struct {
	f         *os.File
	stopEvent windows.Handle
	decoder   utf16x.Decoder
	buffer    bytes.Buffer
	resizeCh  chan os.Signal
}

type bufferResizeSignal struct{}

func (s bufferResizeSignal) String() string { return "buffer resize signal" }
func (s bufferResizeSignal) Signal()        {}

func NewInput(f *os.File) (*Input, error) {
	event, err := windows.CreateEvent(nil, 0, 0, nil)

	if err != nil {
		return nil, err
	}

	r := &Input{
		f:         f,
		stopEvent: event,
		resizeCh:  make(chan os.Signal, 1),
	}

	return r, nil
}

func (i *Input) MakeRaw() (func() error, error) {
	fd := int(i.f.Fd())
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

func (i *Input) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	if i.buffer.Len() > 0 {
		return i.buffer.Read(p)
	}

	console := windows.Handle(i.f.Fd())
	buffer := make([]windowsx.INPUT_RECORD, 1024)

	for {
		n, err := windowsx.PeekConsoleInput(console, buffer)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			if err := i.waitEvent(ctx); err != nil {
				return 0, err
			}
			continue
		}

		n, err = windowsx.ReadConsoleInput(console, buffer)
		if err != nil {
			return 0, err
		}

		for idx := range n {
			record := buffer[idx]

			if _, ok := record.WindowsBufferSizeEvent(); ok {
				select {
				case i.resizeCh <- bufferResizeSignal{}:
				default:
				}
			}

			keyEvent, ok := record.KeyEvent()
			if !ok || keyEvent.KeyDown == 0 {
				continue
			}

			if keyEvent.VirtualKeyCode != 0 && keyEvent.UnicodeChar == 0 {
				continue
			}

			r, ok := i.decoder.Decode(keyEvent.UnicodeChar)
			if !ok {
				continue
			}

			var runeBytes [utf8.UTFMax]byte
			runeLength := utf8.EncodeRune(runeBytes[:], r)
			i.buffer.Write(runeBytes[:runeLength])
		}

		if i.buffer.Len() > 0 {
			return i.buffer.Read(p)
		}
	}
}

func (i *Input) waitEvent(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	console := windows.Handle(i.f.Fd())

	stop := context.AfterFunc(ctx, func() {
		windows.SetEvent(i.stopEvent)
	})
	defer stop()

	handles := []windows.Handle{console, i.stopEvent}
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

func (i *Input) ResizeC() <-chan os.Signal {
	return i.resizeCh
}

func (i *Input) Close() error {
	close(i.resizeCh)

	err := errors.Join(
		windows.SetEvent(i.stopEvent),
		windows.CloseHandle(i.stopEvent))

	return err
}
