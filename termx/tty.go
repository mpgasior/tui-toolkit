package termx

import (
	"errors"
	"os"
)

type TTY struct {
	In  *os.File
	Out *os.File
}

func (tty *TTY) Close() error {
	err := errors.Join(
		tty.In.Close(),
		tty.Out.Close())

	return err
}
