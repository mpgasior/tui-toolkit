package termx

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"golang.org/x/term"
)

type TerminalOutput interface {
	io.Writer
	Flush() error
	Close() error
	GetSize() (w int, h int, err error)
	EnterAltScreen() (restore func() error, err error)
}

type terminalOutput struct {
	f   *os.File
	buf *bufio.Writer
}

func NewTerminalOutput(out *os.File) (*terminalOutput, error) {
	to := &terminalOutput{
		f:   out,
		buf: bufio.NewWriter(out),
	}

	return to, nil
}

func (to *terminalOutput) Write(p []byte) (int, error) {
	return to.buf.Write(p)
}

func (to *terminalOutput) Flush() error {
	return to.buf.Flush()
}

func (to *terminalOutput) Close() error {
	return to.Flush()
}

func (to *terminalOutput) GetSize() (int, int, error) {
	fd := int(to.f.Fd())
	w, h, err := term.GetSize(fd)

	return w, h, err
}

func (to *terminalOutput) EnterAltScreen() (func() error, error) {
	io.WriteString(to.f, "\033[?1049h")

	var once sync.Once
	var err error
	restore := func() error {
		once.Do(func() {
			fmt.Fprint(to.f, "\033[?1049l")
		})

		return err
	}

	return restore, nil
}
