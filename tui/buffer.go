package tui

import (
	"errors"
	"slices"
)

type Cell struct {
	Primary   rune
	Combining []rune
	Style     Style
	Width     uint8
}

func (c *Cell) Equal(other *Cell) bool {
	if c.Primary != other.Primary {
		return false
	}
	if !slices.Equal(c.Combining, other.Combining) {
		return false
	}
	if c.Style != other.Style {
		return false
	}
	if c.Width != other.Width {
		return false
	}

	return true
}

var ErrInvalidPos = errors.New("invalid position")

type Buffer struct {
	w, h  int
	cells []Cell

	cursorX, cursorY int
}

func NewBuffer(w, h int) *Buffer {
	b := &Buffer{
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

func (b *Buffer) Size() (width, heigh int) {
	return b.w, b.h
}

func (b *Buffer) GetAt(x, y int) (*Cell, error) {
	if x < 0 || x >= b.w || y < 0 || y >= b.h {
		return nil, ErrInvalidPos
	}

	idx := (y * b.w) + x
	return &b.cells[idx], nil
}

func (b *Buffer) SetAt(x, y int, primary rune, combs []rune, width uint8, style Style) {
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

func (b *Buffer) SetCursorPos(x, y int) {
	b.cursorX, b.cursorY = x, y
}
