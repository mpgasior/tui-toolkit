package screen

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mpgasior/tui-toolkit/vt"
)

type Screen interface {
	Mutator
	io.WriterTo
	Flush(w io.Writer) (int64, error)
	Resize(w, h int)
}

type screen struct {
	next       Buffer
	current    Buffer
	buf        bytes.Buffer
	styleCache map[Style]string
}

func New(w, h int) Screen {
	return &screen{
		next:       NewBuffer(w, h),
		current:    NewBuffer(w, h),
		styleCache: make(map[Style]string),
	}
}

func (s *screen) Size() (x, y int) {
	return s.next.Size()
}

func (s *screen) Resize(w, h int) {
	oldW, oldH := s.Size()
	if oldW == w && oldH == h {
		return
	}

	s.next = NewBuffer(w, h)
	s.current = NewBuffer(w, h)
}

func (s *screen) SetAt(x, y int, primary rune, combs []rune, width uint8, style Style) {
	s.next.SetAt(x, y, primary, combs, width, style)
}

func (s *screen) SetCursorPos(x, y int) {
	s.next.SetCursorPos(x, y)
}

func (s *screen) WriteTo(writer io.Writer) (n int64, err error) {
	n, err = s.writeTo(writer, false)
	s.current, s.next = s.next, s.current

	return n, err
}

func (s *screen) Flush(writer io.Writer) (n int64, err error) {
	s.current, s.next = s.next, s.current
	n, err = s.writeTo(writer, true)
	s.current, s.next = s.next, s.current

	return n, err
}

func (s *screen) writeTo(writer io.Writer, force bool) (n int64, err error) {
	w, h := s.next.Size()

	s.buf.WriteString(vt.CursorHide)
	for row := range h {
		for col := 0; col < w; {
			current, _ := s.current.GetAt(col, row)
			next, _ := s.next.GetAt(col, row)

			if !force && next.Equal(current) {
				col += int(next.Width)
				continue
			}

			s.moveCursor(col, row)

			if _, ok := s.styleCache[next.Style]; !ok {
				s.styleCache[next.Style] = next.Style.Sequence()
			}

			s.buf.WriteString(s.styleCache[next.Style])
			s.buf.WriteRune(next.Primary)
			if next.Combining != nil {
				s.buf.WriteString(string(next.Combining))
			}
			s.buf.WriteString(vt.SGRReset)

			col += int(next.Width)
		}
	}

	cursorX, cursorY := s.next.GetCursorPos()

	if cursorX != -1 && cursorY != -1 {
		s.moveCursor(cursorX, cursorY)
		s.buf.WriteString(vt.CursorShow)
	}

	s.buf.WriteString(vt.SGRReset)
	return s.buf.WriteTo(writer)
}

func (s *screen) moveCursor(x, y int) {
	fmt.Fprintf(&s.buf, vt.CursorPositionFmt, y+1, x+1)
}
