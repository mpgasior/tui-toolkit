package ui

import (
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/models"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type Search struct {
	*models.TextInput
	*models.Spinner

	searching bool
}

func NewSearch() Search {
	return Search{
		TextInput: &models.TextInput{},
		Spinner:   models.NewSpinner("spinner"),
	}
}

func (s *Search) SetSearching(searching bool) {
	s.searching = searching
}

func (s *Search) Draw(vp view.Port, focused bool) {
	boxStyle := screen.DefaultStyle
	if focused {
		boxStyle = boxStyle.Fg(screen.ColorGreen)
	}
	draw.Box(vp, draw.BoxBorderRounded, boxStyle)

	layout := view.SplitV(vp.Offset(1),
		view.Dynamic("input", 1),
		view.Fixed("spinner", 1),
	)

	inputView, spinnerView := layout["input"], layout["spinner"]

	w, _ := inputView.Size()
	runes, cursor := s.TextInput.Slice(w)
	if len(runes) == 0 {
		draw.Line(inputView, "Search...", screen.DefaultStyle.Fg(screen.ColorBlue))
	} else {
		draw.Line(inputView, string(runes), screen.DefaultStyle)
	}

	if focused {
		inputView.SetCursorPos(cursor, 0)
	}

	if s.searching {
		r := s.Frame()
		draw.Rune(spinnerView, 0, 0, r, screen.DefaultStyle)
	}
}
