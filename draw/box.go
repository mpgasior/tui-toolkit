package draw

import "github.com/mpgasior/tui-toolkit/gfx"

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

type BoxOptions struct {
	Border BoxBorder
	Style  gfx.Style
}

var DefaultBoxOptions = BoxOptions{
	Border: BoxBorderThin,
}

func Box(sf SurfaceMutator, opts BoxOptions) {
	var zero BoxOptions
	if opts == zero {
		opts = DefaultBoxOptions
	}

	border := opts.Border
	style := opts.Style

	w, h := sf.Size()
	if w < 2 || h < 2 {
		return
	}

	for col := range w {
		sf.Set(col, 0, border.Top, style)
		sf.Set(col, h-1, border.Bottom, style)
	}

	for row := range h {
		sf.Set(0, row, border.Left, style)
		sf.Set(w-1, row, border.Right, style)
	}

	sf.Set(0, 0, border.TopLeft, style)
	sf.Set(w-1, 0, border.TopRight, style)
	sf.Set(w-1, h-1, border.BottomRight, style)
	sf.Set(0, h-1, border.BottomLeft, style)
}
