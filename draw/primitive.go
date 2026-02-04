package draw

import (
	"github.com/mpgasior/tui-toolkit/screen"
)

func Rune(b screen.Buffer, x, y int, r rune, style screen.Style) {
	b.SetAt(x, y, r, nil, 1, style)
}

func RuneWide(b screen.Buffer, x, y int, r rune, width uint8, style screen.Style) {
	b.SetAt(x, y, r, nil, width, style)
}
