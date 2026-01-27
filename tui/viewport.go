package tui

type Viewport struct {
	buffer *Buffer
	x, y   int
	w, h   int
}

func NewViewport(buf *Buffer) Viewport {
	w, h := buf.Size()

	return Viewport{
		buffer: buf,
		x:      0, y: 0,
		w: w, h: h,
	}
}

func (vp Viewport) Size() (w, h int) {
	return vp.w, vp.h
}

func (vp Viewport) SetAt(x, y int, primary rune, combs []rune, width uint8, style Style) {
	if x < 0 || x >= vp.w || y < 0 || y >= vp.h {
		return
	}

	vp.buffer.SetAt(x+vp.x, y+vp.y, primary, combs, width, style)
}

func (vp Viewport) SetCursorPos(x, y int) {
	if x < 0 || x >= vp.w || y < 0 || y >= vp.h {
		vp.buffer.SetCursorPos(-1, -1)
		return
	}
	vp.buffer.SetCursorPos(vp.x+x, vp.y+y)
}

func (vp Viewport) Slice(x, y, w, h int) Viewport {
	originX := max(0, min(x, vp.w))
	originY := max(0, min(y, vp.h))

	remainingW := vp.w - originX
	remainingH := vp.h - originY

	return Viewport{
		buffer: vp.buffer,
		x:      vp.x + originX,
		y:      vp.y + originY,
		w:      max(0, min(w, remainingW)),
		h:      max(0, min(h, remainingH)),
	}
}
