package draw

import "github.com/mpgasior/tui-toolkit/gfx"

type Surface struct {
	buf           *gfx.Buffer
	dx, dy        int
	Width, Height int
}

func NewSurface(buf *gfx.Buffer) *Surface {
	return &Surface{
		buf:   buf,
		Width: buf.Width, Height: buf.Height,
	}
}

func (sf *Surface) Size() (w, h int) {
	return sf.Width, sf.Height
}

func (sf *Surface) Set(x, y int, r rune, s gfx.Style) {
	if x < 0 || x >= sf.Width || y < 0 || y >= sf.Height {
		return
	}

	sf.buf.Set(x+sf.dx, y+sf.dy, r, s)
}

func (sf *Surface) SetCursor(x, y int) {
	sf.buf.SetCursor(x+sf.dx, y+sf.dy)
}

func (sf *Surface) Inset(top, right, bottom, left int) *Surface {
	dx, dy := max(0, left), max(0, top)
	dw, dh := left+right, top+bottom

	return &Surface{
		buf:    sf.buf,
		dx:     sf.dx + dx,
		dy:     sf.dy + dy,
		Width:  max(0, sf.Width-dw),
		Height: max(0, sf.Height-dh),
	}
}

func (sf *Surface) Sub(x, y, w, h int) *Surface {
	x, y = max(0, x), max(0, y)
	w, h = max(0, w), max(0, h)

	if x >= sf.Width || y >= sf.Height || x > sf.Width || y > sf.Height {
		return &Surface{
			buf:   sf.buf,
			Width: 0, Height: 0,
		}
	}

	return &Surface{
		buf: sf.buf,
		dx:  x + sf.dx, dy: y + sf.dy,
		Width: w, Height: h,
	}
}
