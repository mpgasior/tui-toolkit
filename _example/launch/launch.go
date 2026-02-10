package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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

		if keyEvent.IsKey(vt.KeyE) {
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

		if keyEvent.IsKey(vt.KeyF) {
			return mvu.TaskOne(mvu.LaunchEvent{
				CmdBuilder: func(ttyIn, ttyOut *os.File) (cmd *exec.Cmd, captureOutput bool, err error) {
					data := bytes.NewBufferString("a\nb\nc\nd\n")

					cmd = exec.Command("fzf")
					cmd.Stdin = data
					cmd.Stdout = ttyOut
					cmd.Stderr = os.Stderr
					return cmd, true, nil
				},
				OnResult: func(out []byte, err error) mvu.Task {
					result, _ := strings.CutSuffix(string(out), "\n")
					fmt.Printf("Fzf: '%s' %v\n", result, err)
					reader := bufio.NewReader(os.Stdin)

					fmt.Printf("%sPress any key to dismiss %s", vt.FormatSGR(vt.FgGreen), vt.SGRReset)
					reader.ReadRune()
					return mvu.TaskNone
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
	text := `
	Press e for nvim
	Press f for fzf
	Press ctrl+c to exit
	`
	draw.Lines(center.Offset(1), text, screen.DefaultStyle)
}

func main() {
	if err := mvu.Run(&App{}); err != nil {
		log.Fatal(err)
	}
}
