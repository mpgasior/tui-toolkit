package ui

import (
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/mvu"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/vt"
)

type SearchInput struct {
	Term []rune
}

func (s *SearchInput) Init() mvu.Task {
	return mvu.TaskNone
}

func (s *SearchInput) Update(e mvu.Event) mvu.Task {
	if key, ok := e.(vt.KeyEvent); ok {
		if key.IsKey(vt.KeyEsc) {
			s.Term = nil
			return mvu.TaskNone
		}

		if key.IsKey(vt.KeyBackspace) {
			if len(s.Term) > 0 {
				s.Term = s.Term[:len(s.Term)-1]
			}
			return mvu.TaskNone
		}

		if key.Rune != utf8.RuneError {
			s.Term = append(s.Term, key.Rune)
		}
	}

	return mvu.TaskNone
}

func (s *SearchInput) Render(ctx mvu.RenderContext) {
	v := ctx.View
	boxStyle := screen.DefaultStyle
	if ctx.Focused {
		boxStyle = boxStyle.Fg(screen.ColorGreen)
	}

	draw.Box(v, draw.BoxBorderThin, boxStyle)

	body := v.Offset(1)
	if len(s.Term) == 0 {
		style := screen.DefaultStyle.Fg(screen.ColorBlue)
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
