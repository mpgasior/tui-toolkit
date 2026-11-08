package termio

import (
	"errors"
	"fmt"

	"golang.org/x/sys/unix"
)

type writer struct {
	fd   int
	epFd int
}

func NewWriter(fd int) (_ *writer, err error) {
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

	epfd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		return nil, err
	}
	defer closeFd(epfd, "epoll fd")

	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{
		Events: unix.EPOLLOUT,
		Fd:     int32(fd),
	}); err != nil {
		return nil, err
	}

	return &writer{
		fd:   fd,
		epFd: epfd,
	}, nil
}

func (w *writer) Close() error {
	return errors.Join(
		unix.SetNonblock(w.fd, false),
		unix.Close(w.epFd),
	)
}

func (w *writer) Write(p []byte) (n int, err error) {
	total := 0

	for total < len(p) {
		n, err = unix.Write(w.fd, p[total:])
		if n > 0 {
			total += n
		}

		if err == nil || err == unix.EINTR {
			continue
		}

		if err == unix.EAGAIN || err == unix.EWOULDBLOCK {
			events := make([]unix.EpollEvent, 1)
			for {
				_, err := unix.EpollWait(w.epFd, events, -1)
				if err != nil {
					if err == unix.EINTR {
						continue
					}
					return total, err
				}
				break
			}
		}
	}

	return total, nil
}
