package termx

import (
	"errors"
	"os"
)

func OpenTTY() (*TTY, error) {
	conin, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	conout, err := os.OpenFile("CONOUT$", os.O_RDWR, 0)

	if err != nil {
		err := errors.Join(
			err, conin.Close())

		return nil, err
	}

	tty := &TTY{
		In:  conin,
		Out: conout,
	}

	return tty, nil
}
