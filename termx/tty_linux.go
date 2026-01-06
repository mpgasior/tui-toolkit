package termx

import (
	"os"
)

func OpenTTY() (*TTY, error) {
	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	tty := &TTY{
		In:  f,
		Out: f,
	}

	return tty, nil
}
