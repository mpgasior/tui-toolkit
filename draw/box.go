package draw

import "github.com/mpgasior/tui-toolkit/screen"

type BoxBorder struct {
	TopLeft, Top, TopRight, Right, BottomRight, Bottom, BottomLeft, Left rune
}

var (
	BoxBorderThin   = BoxBorder{'┌', '─', '┐', '│', '┘', '─', '└', '│'}
	BoxBorderDouble = BoxBorder{'╔', '═', '╗', '║', '╝', '═', '╚', '║'}
)

func Box(b screen.Buffer, border BoxBorder, style screen.Style) {
	w, h := b.Size()
	if w < 2 || h < 2 {
		return
	}

	for col := range w {
		Rune(b, col, 0, border.Top, style)
		Rune(b, col, h-1, border.Bottom, style)
	}

	for row := range h {
		Rune(b, 0, row, border.Left, style)
		Rune(b, w-1, row, border.Right, style)
	}

	Rune(b, 0, 0, border.TopLeft, style)
	Rune(b, w-1, 0, border.TopRight, style)
	Rune(b, w-1, h-1, border.BottomRight, style)
	Rune(b, 0, h-1, border.BottomLeft, style)
}
