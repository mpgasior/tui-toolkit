package gfx

import (
	"io"

	"github.com/mpgasior/tui-toolkit/ansi"
)

func WriteStyleDiff(w io.Writer, prev, next Style) (int, error) {
	var stackBuf [64]byte
	buf := AppendStyleDiff(stackBuf[:], prev, next)

	return w.Write(buf)
}

func WriteBufferDiff(w io.Writer, prev, next *Buffer) (int, error) {
	var buf []byte
	buf = AppendBufferDiff(buf, prev, next)
	return w.Write(buf)
}

func WriteReset(w io.Writer) (int, error) {
	var buf [64]byte
	return w.Write(ansi.AppendSGRReset(buf[:0]))
}
