package screen

import (
	"errors"
)

type Buffer interface {
	Size() (width, height int)
	GetAt(x, y int) (*Cell, error)
	SetAt(x, y int, primary rune, combs []rune, width uint8, style Style)
	SetCursorPos(x, y int)
}

var ErrInvalidPos = errors.New("invalid position")

type cellBuffer struct {
	w, h  int
	cells []Cell

	cursorX, cursorY int
}

func New(w, h int) Buffer {
	b := &cellBuffer{
		w: w, h: h,
		cells: make([]Cell, w*h),

		cursorX: -1,
		cursorY: -1,
	}

	defaultCell := Cell{
		Primary: ' ',
		Width:   1,
		Style:   DefaultStyle,
	}

	for i := range b.cells {
		b.cells[i] = defaultCell
	}

	return b
}

func (b *cellBuffer) Size() (width, heigh int) {
	return b.w, b.h
}

func (b *cellBuffer) GetAt(x, y int) (*Cell, error) {
	if x < 0 || x >= b.w || y < 0 || y >= b.h {
		return nil, ErrInvalidPos
	}

	idx := (y * b.w) + x
	return &b.cells[idx], nil
}

func (b *cellBuffer) SetAt(x, y int, primary rune, combs []rune, width uint8, style Style) {
	if x < 0 || x >= b.w || y < 0 || y >= b.h {
		return
	}

	idx := (y * b.w) + x
	cell := &b.cells[idx]

	cell.Primary = primary
	cell.Combining = combs
	cell.Width = width
	cell.Style = style
}

func (b *cellBuffer) SetCursorPos(x, y int) {
	b.cursorX, b.cursorY = x, y
}
