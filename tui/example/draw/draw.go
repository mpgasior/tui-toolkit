package main

import (
	"log"

	"github.com/mpgasior/tui-go/tui"
	"github.com/mpgasior/tui-go/vt"
)

type Component struct{}

func (c Component) Init() tui.Task { return tui.TaskNone }
func (c Component) Update(e tui.Event) tui.Task {
	if keyEvent, ok := e.(tui.KeyEvent); ok {
		if keyEvent.IsKey(vt.KeyEsc, vt.KeyCtrlC, vt.KeyQ) {
			return tui.TaskShutdown
		}
	}
	return tui.TaskNone
}

func (c Component) Render(ctx tui.RenderContext) {
	style := tui.DefaultStyle.
		Bg(tui.ColorHex(0x008080))

	w, h := ctx.Viewport.Size()

	tui.DrawClear(ctx.Viewport, style)
	tui.DrawBox(ctx.Viewport.Slice(0, 0, w, h-1), tui.BoxBorderThin, style)
	tui.DrawText(ctx.Viewport.Slice(0, h-1, w, 1), "[esc/ctrl+c/q] quit", style.Fg(tui.ColorYellow))
}

func main() {
	if err := tui.Run(Component{}); err != nil {
		log.Fatal(err)
	}
}
