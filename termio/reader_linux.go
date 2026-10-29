package termio

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sys/unix"
)

type reader struct {
	fd       int
	cancelFd int
	epFd     int
}

func NewReader(fd int) (_ *reader, err error) {
	closeFd := func(fd int, info string) {
		if err != nil {
			if closeErr := unix.Close(fd); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("termio: failed to close %s: %w", info, closeErr))
			}
		}
	}

	if err := unix.SetNonblock(fd, true); err != nil {
		return nil, err
	}

	cancelFd, err := unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC)
	if err != nil {
		return nil, fmt.Errorf("termio: failed to create event fd: %w", err)
	}
	defer closeFd(cancelFd, "cancel fd")

	epfd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		return nil, err
	}
	defer closeFd(epfd, "epoll fd")

	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(fd),
	}); err != nil {
		return nil, err
	}

	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, cancelFd, &unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(cancelFd),
	}); err != nil {
		return nil, err
	}

	return &reader{
		fd:       fd,
		epFd:     epfd,
		cancelFd: cancelFd,
	}, nil
}

func (r *reader) Close() error {
	return errors.Join(
		unix.SetNonblock(r.fd, false),
		unix.Close(r.epFd),
		unix.Close(r.cancelFd),
	)
}

func (r *reader) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	stopCancel := context.AfterFunc(ctx, func() {
		buf := [8]byte{1}
		for {
			if _, err := unix.Write(r.cancelFd, buf[:]); err != nil {
				if err == unix.EINTR {
					continue
				}
			}
			break
		}
	})
	defer stopCancel()

	events := make([]unix.EpollEvent, 1)
	for {
		_, err := unix.EpollWait(r.epFd, events, -1)
		if err != nil {
			if errors.Is(err, unix.EINTR) {
				continue
			}
			return 0, err
		}
		break
	}

	switch events[0].Fd {
	case int32(r.fd):
		for {
			n, err := unix.Read(r.fd, p)
			if err != nil {
				if err == unix.EINTR {
					continue
				}
			}

			return n, err
		}
	case int32(r.cancelFd):
		var buf [8]byte
		for {
			if _, err := unix.Read(r.cancelFd, buf[:]); err != nil {
				if err == unix.EINTR {
					continue
				}
				if err == unix.EAGAIN || err == unix.EWOULDBLOCK {
					return 0, ctx.Err()
				}
				return 0, fmt.Errorf("termio: cancel read failed: %w", err)
			}
			break
		}
		return 0, ctx.Err()
	default:
		return 0, errors.New("termio: unknown file descriptor event")
	}
}
