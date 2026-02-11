package draw

import "github.com/mpgasior/tui-toolkit/screen"

type BoxBorder struct {
	TopLeft, Top, TopRight, Right, BottomRight, Bottom, BottomLeft, Left rune
}

var (
	BoxBorderThin    = BoxBorder{'┌', '─', '┐', '│', '┘', '─', '└', '│'}
	BoxBorderDouble  = BoxBorder{'╔', '═', '╗', '║', '╝', '═', '╚', '║'}
	BoxBorderHeavy   = BoxBorder{'┏', '━', '┓', '┃', '┛', '━', '┗', '┃'}
	BoxBorderRounded = BoxBorder{'╭', '─', '╮', '│', '╯', '─', '╰', '│'}
	BoxBorderASCII   = BoxBorder{'+', '-', '+', '|', '+', '-', '+', '|'}
	BoxBorderCorners = BoxBorder{'┌', ' ', '┐', ' ', '┘', ' ', '└', ' '}
)

func Box(m screen.Mutator, border BoxBorder, style screen.Style) {
	w, h := m.Size()
	if w < 2 || h < 2 {
		return
	}

	for col := range w {
		Rune(m, col, 0, border.Top, style)
		Rune(m, col, h-1, border.Bottom, style)
	}

	for row := range h {
		Rune(m, 0, row, border.Left, style)
		Rune(m, w-1, row, border.Right, style)
	}

	Rune(m, 0, 0, border.TopLeft, style)
	Rune(m, w-1, 0, border.TopRight, style)
	Rune(m, w-1, h-1, border.BottomRight, style)
	Rune(m, 0, h-1, border.BottomLeft, style)
}
