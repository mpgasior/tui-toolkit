package gfx

import (
	"io"
)

type Screen struct {
	Front, Back *Buffer
	buf         []byte
}

func NewScreen(w, h int, buf []byte) *Screen {
	front := NewBuffer(w, h)
	back := NewBuffer(w, h)

	return &Screen{
		Front: front, Back: back,
		buf: buf,
	}
}

func (s *Screen) Swap() {
	s.Front, s.Back = s.Back, s.Front
}

func (s *Screen) WriteTo(w io.Writer) (int64, error) {
	s.buf = s.buf[:0]
	s.buf = AppendBufferDiff(s.buf, s.Back, s.Front)

	n, err := w.Write(s.buf)
	return int64(n), err
}

func (s *Screen) Flush(w io.Writer) (int64, error) {
	s.buf = s.buf[:0]
	s.buf = AppendBufferDiff(s.buf, nil, s.Front)

	n, err := w.Write(s.buf)
	return int64(n), err
}
