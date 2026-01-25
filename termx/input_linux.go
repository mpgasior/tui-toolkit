//go:build linux

package termx

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

type terminalInput struct {
	f     *os.File
	epfd  int
	pipeR *os.File
	pipeW *os.File
}

func NewTerminalInput(f *os.File) (TerminalInput, error) {
	var success bool
	epfd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		return nil, err
	}
	defer func() {
		if !success {
			unix.Close(int(f.Fd()))
		}
	}()

	rd, wr, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer func() {
		if !success {
			rd.Close()
			wr.Close()
		}
	}()

	err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(os.Stdin.Fd()), &unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(f.Fd()),
	})
	if err != nil {
		return nil, err
	}

	err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(rd.Fd()), &unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(rd.Fd()),
	})
	if err != nil {
		return nil, err
	}

	inputReader := &terminalInput{
		f:     f,
		epfd:  epfd,
		pipeR: rd,
		pipeW: wr,
	}

	success = true
	return inputReader, nil
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

func (ti *terminalInput) ReadContext(ctx context.Context, p []byte) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	stop := context.AfterFunc(ctx, func() {
		_, _ = ti.pipeW.Write([]byte{0})
	})
	defer stop()

	events := make([]unix.EpollEvent, 1)
	for {
		_, err := unix.EpollWait(ti.epfd, events, -1)
		if err != nil {
			if err == unix.EINTR {
				continue
			}

			return 0, err
		}

		if int(events[0].Fd) == int(ti.f.Fd()) {
			return ti.f.Read(p)
		}

		if int(events[0].Fd) == int(ti.pipeR.Fd()) {
			if _, err := io.CopyN(io.Discard, ti.pipeR, 1); err != nil {
				return 0, err
			}
			return 0, ctx.Err()
		}
	}
}

func (ti *terminalInput) Close() error {
	err := errors.Join(
		ti.pipeW.Close(),
		ti.pipeR.Close(),
		unix.Close(ti.epfd))

	return err
}
