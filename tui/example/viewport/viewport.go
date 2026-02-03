package main

import (
	"log"

	"github.com/mpgasior/tui-go/tui"
	"github.com/mpgasior/tui-go/vt"
)

type Component struct{}

func (c Component) Init() tui.Task { return tui.TaskNone }
func (c Component) Update(e tui.Event) tui.Task {
	if k, ok := e.(tui.KeyEvent); ok {
		if k.IsKey(vt.KeyCtrlC, vt.KeyQ) {
			return tui.TaskShutdown
		}
	}
	return tui.TaskNone
}

func (c Component) Render(ctx tui.RenderContext) {
	main := ctx.Viewport.SplitH(
		tui.ViewportConstraintFixed("header", 3),
		tui.ViewportConstraintDynamic("body", 1),
		tui.ViewportConstraintFixed("footer", 3),
	)

	tui.DrawBox(main["header"], tui.BoxBorderThin, tui.DefaultStyle.Fg(tui.ColorHex(0x303030)))
	tui.DrawBox(main["footer"], tui.BoxBorderThin, tui.DefaultStyle.Fg(tui.ColorHex(0x303030)))

	tui.DrawText(main["header"].Offset(1), "Layouts are easy to use", tui.DefaultStyle.Fg(tui.ColorGreen))
	tui.DrawText(main["footer"].Offset(1), "[ctrl+c/q] exit", tui.DefaultStyle)

	body := main["body"].SplitV(
		tui.ViewportConstraintFixed("left-margin", 1),
		tui.ViewportConstraintDynamic("col-1", 1),
		tui.ViewportConstraintDynamic("col-2", 1),
		tui.ViewportConstraintDynamic("col-3", 1),
		tui.ViewportConstraintFixed("right-margin", 1))

	tui.DrawFill(body["left-margin"], '+', tui.DefaultStyle)
	tui.DrawFill(body["right-margin"], '+', tui.DefaultStyle)
	tui.DrawBox(body["col-1"], tui.BoxBorderDouble, tui.DefaultStyle.Fg(tui.ColorBlue))
	tui.DrawBox(body["col-2"], tui.BoxBorderDouble, tui.DefaultStyle.Fg(tui.ColorYellow))
	tui.DrawBox(body["col-3"], tui.BoxBorderDouble, tui.DefaultStyle.Fg(tui.ColorRed))
}

func main() {
	layout := Component{}

	if err := tui.Run(layout); err != nil {
		log.Fatal(err)
	}
}
