package gfx

import (
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/ansi"
)

func AppendStyleDiff(dst []byte, prev, next Style) []byte {
	if prev == next {
		return dst
	}

	dst = ansi.AppendSGRStart(dst)

	prevFg, prevBg, prevAttr := prev.Unpack()
	nextFg, nextBg, nextAttr := next.Unpack()

	if prevFg != nextFg {
		dst = nextFg.AppendFgTo(dst)
	}

	if prevBg != nextBg {
		dst = nextBg.AppendBgTo(dst)
	}

	if prevAttr != nextAttr {
		dst = nextAttr.AppendTo(dst)
	}

	return ansi.AppendSGREnd(dst)
}

func AppendBufferDiff(dst []byte, prev, next *Buffer) []byte {
	if prev == nil {
		dst = ansi.AppendSGRReset(dst)
		for y := 0; y < next.Height; y += 1 {
			rowOffset := y * next.Width
			dst = ansi.AppendCursorMove(dst, 0, y)
			for x := 0; x < next.Width; x += 1 {
				pos := rowOffset + x
				cell := next.Cells[pos]

				dst = AppendStyleDiff(dst, StyleDefault, cell.Style)
				dst = utf8.AppendRune(dst, cell.R)
			}
		}
		return dst
	}

	lastX, lastY := -1, -1
	for y := 0; y < next.Height; y += 1 {
		rowOffset := y * next.Width
		for x := 0; x < next.Width; x += 1 {
			pos := rowOffset + x
			prevCell, nextCell := prev.Cells[pos], next.Cells[pos]
			if prevCell == nextCell {
				continue
			}

			prevStyle := StyleDefault
			if x > 0 {
				prevStyle = next.Cells[pos-1].Style
			}

			if x == lastX+1 || y != lastY {
				dst = ansi.AppendCursorMove(dst, x, y)
			}

			dst = AppendStyleDiff(dst, prevStyle, nextCell.Style)
			dst = utf8.AppendRune(dst, nextCell.R)

			lastX, lastY = x, y
		}
	}

	return dst
}
