package tui

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mpgasior/tui-go/vt"
)

type Renderer struct {
	Front *Buffer
	Back  *Buffer
	buf   bytes.Buffer
}

func NewRenderer(w, h int) Renderer {
	return Renderer{
		Front: NewBuffer(w, h),
		Back:  NewBuffer(w, h),
	}
}

func (r *Renderer) Size() (w, h int) {
	return r.Front.Size()
}

func (r *Renderer) SwapBuffers() {
	r.Front, r.Back = r.Back, r.Front
}

func (r *Renderer) WriteTo(writer io.Writer) (n int64, err error) {
	r.buf.WriteString(vt.CursorHide)
	w, h := r.Size()
	for row := range h {
		for col := 0; col < w; {
			front, _ := r.Front.GetAt(col, row)
			back, _ := r.Back.GetAt(col, row)

			if front.Equal(back) {
				col += int(front.Width)
				continue
			}

			r.moveCursor(col, row)

			r.buf.WriteString(front.Style.Sequence())
			r.buf.WriteRune(front.Primary)
			if front.Combining != nil {
				r.buf.WriteString(string(front.Combining))
			}
			r.buf.WriteString(vt.SGRReset)

			col += int(front.Width)
		}
	}

	cursorX, cursorY := r.Front.cursorX, r.Front.cursorY

	if cursorX != -1 && cursorY != -1 {
		r.moveCursor(cursorX, cursorY)
		r.buf.WriteString(vt.CursorShow)
	}

	return r.buf.WriteTo(writer)
}

func (r *Renderer) ForceRedraw() {
	w, h := r.Front.Size()
	for col := range w {
		for row := range h {
			r.Front.SetAt(col, row, ' ', nil, 1, DefaultStyle)
		}
	}
}

func (r *Renderer) moveCursor(x, y int) {
	fmt.Fprintf(&r.buf, vt.CursorPositionFmt, y+1, x+1)
}
