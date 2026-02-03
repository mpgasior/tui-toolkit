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

func (vp Viewport) Offset(offsets ...int) Viewport {
	var top, right, bottom, left int

	switch len(offsets) {
	case 1:
		top, right, bottom, left = offsets[0], offsets[0], offsets[0], offsets[0]
	case 2:
		top, bottom = offsets[0], offsets[0]
		right, left = offsets[1], offsets[1]
	case 3:
		top = offsets[0]
		right, left = offsets[1], offsets[1]
		bottom = offsets[2]
	case 4:
		top, right, bottom, left = offsets[0], offsets[1], offsets[2], offsets[3]
	}

	newX := vp.x + left
	newY := vp.y + top

	newW := max(0, vp.w-left-right)
	newH := max(0, vp.h-top-bottom)

	return Viewport{
		buffer: vp.buffer,
		x:      newX,
		y:      newY,
		w:      newW,
		h:      newH,
	}
}
