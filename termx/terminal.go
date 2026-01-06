package termx

import (
	"errors"
	"os"
)

type Terminal struct {
	*Input
	*Output
}

func New(in *os.File, out *os.File) (*Terminal, error) {
	devIn, err := NewInput(in)
	if err != nil {
		return nil, err
	}

	devOut, err := NewOutput(out)
	if err != nil {
		closeErr := devIn.Close()

		err := errors.Join(closeErr, err)
		return nil, err
	}

	terminal := &Terminal{
		devIn,
		devOut,
	}

	return terminal, nil
}
