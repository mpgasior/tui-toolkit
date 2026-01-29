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

func (r Renderer) Size() (w, h int) {
	return r.Front.Size()
}

func (r Renderer) Swap() {
	r.Front, r.Back = r.Back, r.Front
}

func (r Renderer) Draw(writer io.Writer) {
	w, h := r.Size()
	for row := range h {
		for col := 0; col < w; {
			front, _ := r.Front.GetAt(row, col)
			_, _ = r.Back.GetAt(row, col)

			moveCursor(writer, row, col)

			io.WriteString(writer, front.Style.Sequence())
			io.WriteString(writer, string(front.Primary))
			if front.Combining != nil {
				io.WriteString(writer, string(front.Combining))
			}
			io.WriteString(writer, vt.SGRReset)

			col += int(front.Width)
		}
	}
}

func showCursor(writer io.Writer, show bool) {
	if show {
		io.WriteString(writer, vt.CursorShow)
		return
	}

	io.WriteString(writer, vt.CursorHide)
}

func moveCursor(writer io.Writer, x, y int) {
	fmt.Fprintf(writer, vt.CursorPositionFmt, y+1, x+1)
}
