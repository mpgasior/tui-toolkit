package termx

import (
	"context"
	"sync"

	"golang.org/x/term"
)

type TerminalInput interface {
	MakeRaw() (restore func() error, err error)
	Close() error
	ReadContext(ctx context.Context, p []byte) (n int, err error)
}

func makeRaw(fd int) (func() error, error) {
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
