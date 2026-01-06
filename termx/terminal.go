package termx

import (
	"errors"
	"os"
)

type Terminal interface {
	TerminalInput
	TerminalOutput
}

type terminal struct {
	TerminalInput
	TerminalOutput
}

func NewTerminal(in *os.File, out *os.File) (Terminal, error) {
	termInput, err := NewTerminalInput(in)
	if err != nil {
		return nil, err
	}

	termOutput, err := NewTerminalOutput(out)
	if err != nil {
		closeErr := termInput.Close()

		err := errors.Join(closeErr, err)
		return nil, err
	}

	terminal := &terminal{
		termInput,
		termOutput,
	}

	return terminal, nil
}

func (t *terminal) Close() error {
	inErr := t.TerminalInput.Close()
	outErr := t.TerminalOutput.Close()

	err := errors.Join(inErr, outErr)
	return err
}
