package termx

import (
	"os"

	"golang.org/x/term"
)

type Output struct {
	f *os.File
}

func NewOutput(out *os.File) (*Output, error) {
	to := &Output{
		f: out,
	}

	return to, nil
}

func (o *Output) Write(p []byte) (int, error) {
	return o.f.Write(p)
}

func (o *Output) GetSize() (w int, h int, err error) {
	fd := int(o.f.Fd())
	w, h, err = term.GetSize(fd)

	return w, h, err
}
