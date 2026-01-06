package draw

import (
	"github.com/mpgasior/tui-toolkit/screen"
)

func Rune(m screen.Mutator, x, y int, r rune, style screen.Style) {
	m.SetAt(x, y, r, nil, 1, style)
}

func RuneWide(m screen.Mutator, x, y int, r rune, width uint8, style screen.Style) {
	m.SetAt(x, y, r, nil, width, style)
}
