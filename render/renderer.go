package render

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/vt"
)

type Renderer struct {
	back screen.Buffer
	buf  bytes.Buffer
}

func New(w, h int) Renderer {
	return Renderer{
		back: screen.New(w, h),
	}
}

func (r *Renderer) Render(front screen.Buffer, writer io.Writer) (n int64, err error) {
	fw, fh := front.Size()
	bw, bh := r.back.Size()

	if fw != bw || fh != bh {
		r.back = screen.New(fw, fh)
	}

	r.buf.WriteString(vt.CursorHide)
	for row := range fh {
		for col := 0; col < fw; {
			fCell, _ := front.GetAt(col, row)
			fBack, _ := r.back.GetAt(col, row)

			if fCell.Equal(fBack) {
				col += int(fCell.Width)
				continue
			}

			r.moveCursor(col, row)

			r.buf.WriteString(fCell.Style.Sequence())
			r.buf.WriteRune(fCell.Primary)
			if fCell.Combining != nil {
				r.buf.WriteString(string(fCell.Combining))
			}
			r.buf.WriteString(vt.SGRReset)

			r.back.SetAt(col, row, fCell.Primary, fCell.Combining, fCell.Width, fCell.Style)
			col += int(fCell.Width)
		}
	}

	cursorX, cursorY := front.GetCursorPos()

	if cursorX != -1 && cursorY != -1 {
		r.moveCursor(cursorX, cursorY)
		r.buf.WriteString(vt.CursorShow)
	}

	return r.buf.WriteTo(writer)
}

func (r *Renderer) moveCursor(x, y int) {
	fmt.Fprintf(&r.buf, vt.CursorPositionFmt, y+1, x+1)
}
