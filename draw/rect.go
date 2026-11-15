package draw

import "github.com/mpgasior/tui-toolkit/gfx"

func Fill(m SurfaceMutator, r rune, s gfx.Style) {
	w, h := m.Size()
	for col := range w {
		for row := range h {
			m.Set(col, row, r, s)
		}
	}
}

func Clear(m SurfaceMutator) {
	Fill(m, ' ', gfx.StyleDefault)
}
