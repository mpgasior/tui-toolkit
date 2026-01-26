package tui

type View struct {
	buffer *Buffer
	x, y   int
	w, h   int
}

func NewView(buf *Buffer) View {
	w, h := buf.Size()

	return View{
		buffer: buf,
		x:      0, y: 0,
		w: w, h: h,
	}
}

func (v View) Size() (w, h int) {
	return v.w, v.h
}

func (v View) SetAt(x, y int, primary rune, combs []rune, width uint8, style Style) {
	if x < 0 || x >= v.w || y < 0 || y >= v.h {
		return
	}

	v.buffer.SetAt(x+v.x, y+v.y, primary, combs, width, style)
}

func (v View) SetRuneAt(x, y int, r rune, style Style) {
	v.SetAt(x, y, r, nil, 1, style)
}

func (v View) ShowCursor(x, y int) {
	v.buffer.SetCursorPos(v.x+x, v.y+y)
}

func (v View) HideCursor() {
	v.buffer.SetCursorPos(-1, -1)
}

func (v View) SubView(x, y, w, h int) View {
	newX := v.x + x
	newY := v.y + y

	actualW := w
	if x+w > v.w {
		actualW = v.w - x
	}
	actualH := h
	if y+h > v.h {
		actualH = v.h - y
	}

	return View{
		buffer: v.buffer,
		x:      newX,
		y:      newY,
		w:      max(0, actualW),
		h:      max(0, actualH),
	}
}

func (v View) Clear() {
	emptyStyle := NewStyle()
	for row := range v.h {
		for col := range v.w {
			v.SetAt(col, row, ' ', nil, 1, emptyStyle)
		}
	}
}
