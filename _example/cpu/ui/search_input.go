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
	boxStyle := screen.DefaultStyle
	if ctx.Focused {
		boxStyle = boxStyle.Fg(screen.ColorGreen)
	}
	draw.Box(ctx.View, draw.BoxBorderThin, boxStyle)

	body := ctx.View.Offset(1)
	w, _ := body.Size()
	cursorXPos := 0
	if len(s.Term) == 0 {
		style := screen.DefaultStyle.Fg(screen.ColorBlue)
		draw.Line(body, "Search...", style)
	} else if len(s.Term) >= w {
		start := len(s.Term) - w + 1
		runes := s.Term[start:len(s.Term)]

		draw.Line(body, string(runes), screen.DefaultStyle)
		cursorXPos = len(runes)
	} else {
		draw.Line(body, string(s.Term), screen.DefaultStyle)
		cursorXPos = len(s.Term)
	}

	if ctx.Focused {
		body.SetCursorPos(cursorXPos, 0)
	}
}
