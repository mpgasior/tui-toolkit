package gfx

type Cell struct {
	R     rune
	Style Style
}

type Buffer struct {
	Width, Height int
	Cells         []Cell

	CursorX, CursorY int
	CursorStyle      CursorStyle
}

func NewBuffer(w, h int) *Buffer {
	w = max(0, w)
	h = max(0, h)

	return &Buffer{
		Width: w, Height: h,
		Cells: make([]Cell, w*h),
	}
}

func (b *Buffer) Size() (w, h int) {
	return b.Width, b.Height
}

func (b *Buffer) Set(x, y int, r rune, s Style) {
	if x < 0 || x >= b.Width || y < 0 || y >= b.Height {
		return
	}
	b.Cells[y*b.Width+x] = Cell{R: r, Style: s}
}

func (b *Buffer) SetCursor(x, y int, s CursorStyle) {
	if x < 0 || x >= b.Width || y < 0 || y > b.Height {
		b.CursorX, b.CursorY = -1, -1
		return
	}

	b.CursorX, b.CursorY = x, y
	b.CursorStyle = s
}
