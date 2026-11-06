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
	if prev.Width != next.Width || prev.Height != next.Height {
		for y := 0; y < next.Height; y += 1 {
			for x := 0; x < next.Width; x += 1 {
				pos := y*next.Height + x
				cell := next.Cells[pos]
				dst = AppendStyleDiff(dst, StyleDefault, cell.Style)
				dst = utf8.AppendRune(dst, cell.R)
			}
		}
		return dst
	}

	for y := 0; y < next.Height; y += 1 {
		rowOffset := y * next.Width
		for x := 0; x < next.Width; x += 1 {
			pos := rowOffset + x
			prevCell, nextCell := prev.Cells[pos], next.Cells[pos]
			if prevCell == nextCell {
				continue
			}

			dst = ansi.AppendCursorMove(dst, x, y)
			dst = AppendStyleDiff(dst, prevCell.Style, nextCell.Style)
			dst = utf8.AppendRune(dst, nextCell.R)
		}
	}

	return dst
}
