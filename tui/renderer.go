package tui

import (
	"fmt"
	"io"

	"github.com/mpgasior/tui-go/vt"
)

type Renderer struct {
	Front *Buffer
	Back  *Buffer
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

func (r *Renderer) WriteTo(writer io.Writer) {
	r.showCursor(writer, false)
	w, h := r.Size()
	for row := range h {
		for col := 0; col < w; {
			front, _ := r.Front.GetAt(col, row)
			_, _ = r.Back.GetAt(col, row)

			r.moveCursor(writer, col, row)

			io.WriteString(writer, front.Style.Sequence())
			io.WriteString(writer, string(front.Primary))
			if front.Combining != nil {
				io.WriteString(writer, string(front.Combining))
			}
			io.WriteString(writer, vt.SGRReset)

			col += int(front.Width)
		}
	}

	cursorX, cursorY := r.Front.cursorX, r.Front.cursorY

	if cursorX != -1 && cursorY != -1 {
		r.moveCursor(writer, cursorX, cursorY)
		r.showCursor(writer, true)
	}
}

func (r *Renderer) showCursor(writer io.Writer, show bool) {
	if show {
		io.WriteString(writer, vt.CursorShow)
		return
	}

	io.WriteString(writer, vt.CursorHide)
}

func (r *Renderer) moveCursor(writer io.Writer, x, y int) {
	fmt.Fprintf(writer, vt.CursorPositionFmt, y+1, x+1)
}
