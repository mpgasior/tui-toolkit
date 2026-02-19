package ui

import (
	"strings"

	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

func DrawHelp(vp view.Port, focus Focus) {
	items := []string{"Quit: ctrl+c"}

	switch focus {
	case FocusSearch, FocusTable:
		items = append(items, "Focus: [S-]Tab", "Pause: ctrl+p")

		if focus == FocusTable {
			items = append(items,
				"Sort: s",
				"Column: h/l",
				"Top/End: g/G",
			)
		}
	}

	text := strings.Join(items, " • ")

	draw.Line(vp, text, screen.DefaultStyle.Fg(screen.ColorBlue))
}
