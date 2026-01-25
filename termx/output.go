package termx

import (
	"io"
	"os"

	"golang.org/x/term"
)

type TerminalOutput interface {
	io.Writer
	GetSize() (w int, h int, err error)
}

type terminalOutput struct {
	f *os.File
}

func NewTerminalOutput(out *os.File) (*terminalOutput, error) {
	to := &terminalOutput{
		f: out,
	}

	return to, nil
}

func (to *terminalOutput) Write(p []byte) (int, error) {
	return to.f.Write(p)
}

func (to *terminalOutput) GetSize() (int, int, error) {
	fd := int(to.f.Fd())
	w, h, err := term.GetSize(fd)

	return w, h, err
}
