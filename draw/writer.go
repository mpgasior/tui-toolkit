package draw

import (
	"errors"
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/gfx"
)

var ErrInvalidRune = errors.New("draw: invalid rune")

type TextWriter struct {
	sf SurfaceMutator

	row, col int
	style    gfx.Style
	wrap     bool
}

type TextWriterOption = func(tw *TextWriter)

var (
	WithTextWriterStyle = func(s gfx.Style) TextWriterOption {
		return func(tw *TextWriter) {
			tw.style = s
		}
	}

	WithTextWriterWordWrapping = func(wrap bool) TextWriterOption {
		return func(tw *TextWriter) {
			tw.wrap = wrap
		}
	}
)

func NewTextWriter(sf SurfaceMutator, opts ...TextWriterOption) *TextWriter {
	tw := &TextWriter{sf: sf}

	for _, opt := range opts {
		opt(tw)
	}

	return tw
}

func (w *TextWriter) WriteRune(r rune) {
	switch r {
	case '\r':
		w.col = 0
	case '\n':
		w.row += 1
	default:
		w.sf.Set(w.col, w.row, r, w.style)
		w.col += 1
		width, _ := w.sf.Size()
		if w.col >= width && w.wrap {
			w.col, w.row = 0, w.row+1
		}
	}
}

func (w *TextWriter) Write(p []byte) (n int, err error) {
	offset := 0
	for len(p[offset:]) > 0 {
		r, size := utf8.DecodeRune(p[offset:])
		if r == utf8.RuneError && size == 0 {
			return offset, ErrInvalidRune
		}

		w.WriteRune(r)
		offset += size
	}

	return offset, nil
}
