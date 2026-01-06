package termx

import (
	"context"
	"errors"
	"io"
	"os"
	"os/signal"
	"sync"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

type Input struct {
	f        *os.File
	epfd     int
	pipeR    *os.File
	pipeW    *os.File
	resizeCh chan os.Signal
}

func NewInput(f *os.File) (*Input, error) {
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

	inputReader := &Input{
		f:        f,
		epfd:     epfd,
		pipeR:    rd,
		pipeW:    wr,
		resizeCh: make(chan os.Signal, 1),
	}

	signal.Notify(inputReader.resizeCh, unix.SIGWINCH)

	success = true
	return inputReader, nil
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

func (i *Input) ReadContext(ctx context.Context, p []byte) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	stop := context.AfterFunc(ctx, func() {
		_, _ = i.pipeW.Write([]byte{0})
	})
	defer stop()

	events := make([]unix.EpollEvent, 1)
	for {
		_, err := unix.EpollWait(i.epfd, events, -1)
		if err != nil {
			if err == unix.EINTR {
				continue
			}

			return 0, err
		}

		if int(events[0].Fd) == int(i.f.Fd()) {
			return i.f.Read(p)
		}

		if int(events[0].Fd) == int(i.pipeR.Fd()) {
			if _, err := io.CopyN(io.Discard, i.pipeR, 1); err != nil {
				return 0, err
			}
			return 0, ctx.Err()
		}
	}
}

func (i *Input) ResizeC() <-chan os.Signal {
	return i.resizeCh
}

func (i *Input) Close() error {
	close(i.resizeCh)
	err := errors.Join(
		i.pipeW.Close(),
		i.pipeR.Close(),
		unix.Close(i.epfd))

	return err
}
