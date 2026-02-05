package draw

import "github.com/mpgasior/tui-toolkit/screen"

func Fill(m screen.Mutator, r rune, style screen.Style) {
	w, h := m.Size()
	for col := range w {
		for row := range h {
			Rune(m, col, row, r, style)
		}
	}
}

func Clear(m screen.Mutator, style screen.Style) {
	Fill(m, ' ', style)
}
