package gfx

import (
	"io"

	"github.com/mpgasior/tui-toolkit/ansi"
)

func AppendStyleDiff(dst []byte, prev, next Style) []byte {
	if prev == next {
		return dst
	}

	dst = ansi.AppendSGRStart(dst)

	prevFg, prevBg, prevAttr := prev.Unpack()
	nextFg, nextBg, nextAttr := next.Unpack()

	if prevFg != nextFg {
		dst = nextFg.AppendFgTo(dst)
	}

	if prevBg != nextBg {
		dst = nextBg.AppendBgTo(dst)
	}

	if prevAttr != nextAttr {
		dst = nextAttr.AppendTo(dst)
	}

	return ansi.AppendSGREnd(dst)
}

func WriteStyleDiff(w io.Writer, prev, next Style) (int, error) {
	var stackBuf [64]byte
	buf := AppendStyleDiff(stackBuf[:], prev, next)

	return w.Write(buf)
}

func WriteReset(w io.Writer) (int, error) {
	var buf [64]byte
	return w.Write(ansi.AppendSGRReset(buf[:0]))
}
