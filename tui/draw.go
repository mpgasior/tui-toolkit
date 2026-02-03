package tui

import (
	"bufio"
	"strings"
)

func SetRune(vp Viewport, x, y int, r rune, style Style) {
	vp.SetAt(x, y, r, nil, 1, style)
}

func SetRuneWide(vp Viewport, x, y int, r rune, width uint8, style Style) {
	vp.SetAt(x, y, r, nil, width, style)
}

func DrawText(vp Viewport, text string, style Style) {
	scanner := bufio.NewScanner(strings.NewReader(text))

	row := 0
	for scanner.Scan() {
		line := scanner.Text()

		for idx, r := range []rune(line) {
			SetRune(vp, idx, row, r, style)
		}
		row += 1
	}
}

func DrawFill(vp Viewport, r rune, style Style) {
	w, h := vp.Size()
	for col := range w {
		for row := range h {
			SetRune(vp, col, row, r, style)
		}
	}
}

func DrawClear(vp Viewport, style Style) {
	DrawFill(vp, ' ', style)
}

type BoxBorderStyle struct {
	TopLeft, Top, TopRight, Right, BottomRight, Bottom, BottomLeft, Left rune
}

var (
	BoxBorderThin   = BoxBorderStyle{'┌', '─', '┐', '│', '┘', '─', '└', '│'}
	BoxBorderDouble = BoxBorderStyle{'╔', '═', '╗', '║', '╝', '═', '╚', '║'}
)

func DrawBox(vp Viewport, borderStyle BoxBorderStyle, style Style) {
	w, h := vp.Size()
	if w < 2 || h < 2 {
		return
	}

	for col := range w {
		SetRune(vp, col, 0, borderStyle.Top, style)
		SetRune(vp, col, h-1, borderStyle.Bottom, style)
	}

	for row := range h {
		SetRune(vp, 0, row, borderStyle.Left, style)
		SetRune(vp, w-1, row, borderStyle.Right, style)
	}

	SetRune(vp, 0, 0, borderStyle.TopLeft, style)
	SetRune(vp, w-1, 0, borderStyle.TopRight, style)
	SetRune(vp, w-1, h-1, borderStyle.BottomRight, style)
	SetRune(vp, 0, h-1, borderStyle.BottomLeft, style)
}
