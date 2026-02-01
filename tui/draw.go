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

func Fill(vp Viewport, r rune, style Style) {
	w, h := vp.Size()
	for col := range w {
		for row := range h {
			SetRune(vp, col, row, r, style)
		}
	}
}

func Clear(vp Viewport, style Style) {
	Fill(vp, ' ', style)
}

func DrawBox(vp Viewport, border string, style Style) {
	w, h := vp.Size()
	if w < 2 || h < 2 {
		return
	}

	if len(border) < 8 {
		border = "┌─┐│┘─└│"
	}

	b := []rune(border)

	for col := range w {
		SetRune(vp, col, 0, b[1], style)
		SetRune(vp, col, h-1, b[5], style)
	}

	for row := range h {
		SetRune(vp, 0, row, b[7], style)
		SetRune(vp, w-1, row, b[3], style)
	}

	SetRune(vp, 0, 0, b[0], style)
	SetRune(vp, w-1, 0, b[2], style)
	SetRune(vp, w-1, h-1, b[4], style)
	SetRune(vp, 0, h-1, b[6], style)
}
