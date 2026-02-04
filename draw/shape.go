package draw

import "github.com/mpgasior/tui-toolkit/screen"

func Fill(b screen.Buffer, r rune, style screen.Style) {
	w, h := b.Size()
	for col := range w {
		for row := range h {
			Rune(b, col, row, r, style)
		}
	}
}

func Clear(b screen.Buffer, style screen.Style) {
	Fill(b, ' ', style)
}
