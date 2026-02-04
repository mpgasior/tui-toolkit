package view

import (
	"github.com/mpgasior/tui-go/screen"
)

type Port struct {
	buffer screen.Buffer
	x, y   int
	w, h   int
}

func NewPort(buf screen.Buffer) Port {
	w, h := buf.Size()

	return Port{
		buffer: buf,
		x:      0, y: 0,
		w: w, h: h,
	}
}

func (p Port) Size() (w, h int) {
	return p.w, p.h
}

func (p Port) SetAt(x, y int, primary rune, combs []rune, width uint8, style screen.Style) {
	if x < 0 || x >= p.w || y < 0 || y >= p.h {
		return
	}

	p.buffer.SetAt(x+p.x, y+p.y, primary, combs, width, style)
}

func (p Port) SetCursorPos(x, y int) {
	if x < 0 || x >= p.w || y < 0 || y >= p.h {
		p.buffer.SetCursorPos(-1, -1)
		return
	}
	p.buffer.SetCursorPos(p.x+x, p.y+y)
}

func (p Port) Slice(x, y, w, h int) Port {
	originX := max(0, min(x, p.w))
	originY := max(0, min(y, p.h))

	remainingW := p.w - originX
	remainingH := p.h - originY

	return Port{
		buffer: p.buffer,
		x:      p.x + originX,
		y:      p.y + originY,
		w:      max(0, min(w, remainingW)),
		h:      max(0, min(h, remainingH)),
	}
}

func (p Port) Offset(offsets ...int) Port {
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

	newX := p.x + left
	newY := p.y + top

	newW := max(0, p.w-left-right)
	newH := max(0, p.h-top-bottom)

	return Port{
		buffer: p.buffer,
		x:      newX,
		y:      newY,
		w:      newW,
		h:      newH,
	}
}
