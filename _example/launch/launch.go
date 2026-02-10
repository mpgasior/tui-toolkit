package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/mvu"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
	"github.com/mpgasior/tui-toolkit/vt"
)

type App struct{}

func (a *App) Init() mvu.Task { return mvu.TaskNone }

func (a *App) Update(e mvu.Event) mvu.Task {
	if keyEvent, ok := e.(vt.KeyEvent); ok {
		if keyEvent.IsKey(vt.KeyCtrlC) {
			return mvu.TaskShutdown
		}

		if keyEvent.IsKey(vt.KeyE, vt.KeyO) {
			return mvu.TaskOne(mvu.LaunchEvent{
				CmdBuilder: func(ttyIn, ttyOut *os.File) (cmd *exec.Cmd, captureOutput bool, err error) {
					cmd = exec.Command("nvim")
					cmd.Stdin = ttyIn
					cmd.Stdout = ttyOut
					cmd.Stderr = os.Stderr
					return cmd, false, nil
				},
			})
		}
	}

	return mvu.TaskNone
}

func (a *App) Render(ctx mvu.RenderContext) {
	draw.Clear(ctx.View, screen.DefaultStyle)

	center := view.Center(ctx.View, view.Dynamic("l", 1), view.Dynamic("r", 1))
	draw.Box(center, draw.BoxBorderDouble, screen.DefaultStyle)
	draw.Line(center.Offset(1), "Press n for nvim", screen.DefaultStyle)
}

func main() {
	if err := mvu.Run(&App{}); err != nil {
		log.Fatal(err)
	}
}
