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

func AppendCursorStyleDiff(dst []byte, prev, next CursorStyle) []byte {
	if prev == next {
		return dst
	}

	if next == CursorStyleDefault {
		dst = ansi.AppendResetCursorColorRGB(dst)
		return ansi.AppendCursorStyle(dst, ansi.CursorDefault)
	}

	color, shape, blink := next.Unpack()
	cs := ansi.CursorStyle(shape)
	if !blink {
		cs += 1
	}
	dst = ansi.AppendCursorStyle(dst, cs)

	switch color.Kind() {
	case ColorKindRGB:
		hex := color.RGB()
		dst = ansi.AppendCursorColorRGB(dst, hex)
	default:
		dst = ansi.AppendResetCursorColorRGB(dst)
	}

	return dst
}

func AppendBufferDiff(dst []byte, prev, next *Buffer) []byte {
	if prev == nil {
		dst = ansi.AppendSGRReset(dst)
		dst = ansi.AppendCursorHide(dst)
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
		dst = ansi.AppendCursorMove(dst, next.CursorX, next.CursorY)
		dst = AppendCursorStyleDiff(dst, CursorStyleDefault, next.CursorStyle)
		dst = ansi.AppendCursorShow(dst)
		return dst
	}

	dst = ansi.AppendCursorHide(dst)
	dst = ansi.AppendSGRReset(dst)

	curStyle := StyleDefault
	lastX, lastY := -1, -1
	for y := 0; y < next.Height; y += 1 {
		rowOffset := y * next.Width
		for x := 0; x < next.Width; x += 1 {
			pos := rowOffset + x
			prevCell, nextCell := prev.Cells[pos], next.Cells[pos]
			if prevCell == nextCell {
				continue
			}

			if x != lastX+1 || y != lastY {
				dst = ansi.AppendCursorMove(dst, x, y)
				dst = ansi.AppendSGRReset(dst)
				curStyle = StyleDefault
			}

			dst = AppendStyleDiff(dst, curStyle, nextCell.Style)
			dst = utf8.AppendRune(dst, nextCell.R)
			curStyle = nextCell.Style

			lastX, lastY = x, y
		}
	}

	dst = ansi.AppendSGRReset(dst)
	dst = ansi.AppendCursorMove(dst, next.CursorX, next.CursorY)
	dst = AppendCursorStyleDiff(dst, prev.CursorStyle, next.CursorStyle)
	dst = ansi.AppendCursorShow(dst)

	return dst
}
