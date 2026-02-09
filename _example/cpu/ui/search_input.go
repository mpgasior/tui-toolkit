package ui

import (
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
	"github.com/mpgasior/tui-toolkit/vt"
)

type SearchInput struct {
	Term []rune
}

func (s *SearchInput) Update(e vt.KeyEvent) {
	if e.IsKey(vt.KeyBackspace) {
		if len(s.Term) > 0 {
			s.Term = s.Term[:len(s.Term)-1]
		}
		return
	}

	if e.Rune != utf8.RuneError {
		s.Term = append(s.Term, e.Rune)
	}
}

func (s *SearchInput) Draw(v view.Port, focused bool) {
	boxStyle := screen.DefaultStyle
	if focused {
		boxStyle = boxStyle.Fg(screen.ColorGreen)
	}

	draw.Box(v, draw.BoxBorderThin, boxStyle)

	body := v.Offset(1)
	if len(s.Term) == 0 {
		style := screen.DefaultStyle.
			Fg(screen.ColorHex(0x0F0F0F))
		draw.Line(body, "Search...", style)
		body.SetCursorPos(-1, -1)
		return
	}

	w, _ := body.Size()
	if len(s.Term) >= w {
		start := len(s.Term) - w + 1
		runes := s.Term[start:len(s.Term)]

		draw.Line(body, string(runes), screen.DefaultStyle)
		body.SetCursorPos(len(runes), 0)
		return
	}

	draw.Line(body, string(s.Term), screen.DefaultStyle)
	body.SetCursorPos(len(s.Term), 0)
}
