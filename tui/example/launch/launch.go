package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mpgasior/tui-go/tui"
	"github.com/mpgasior/tui-go/vt"
)

type Launcher struct {
	NvimOutput string
	FzfOutput  string
}

func (l *Launcher) Init() tui.Task {
	return tui.TaskNone()
}

func (l *Launcher) Update(e tui.Event) tui.Task {
	if e, ok := e.(tui.KeyEvent); ok {
		if e.IsKey(vt.KeyCtrlC, vt.KeyEsc) {
			return tui.TaskShutdown
		}

		if e.IsRune('n') {
			l.NvimOutput = ""
			launch := tui.LaunchEvent{
				CmdBuilder: func() (cmd *exec.Cmd, captureOutput bool, err error) {
					return exec.Command("nvim"), false, nil
				}, OnResult: func(stdout []byte, err error) tui.Task {
					txt := string(stdout)
					l.NvimOutput = fmt.Sprintf("nvim: '%s' err=%v", txt, err)

					return tui.TaskNone()
				},
			}

			return tui.TaskOne(launch)
		} else if e.IsRune('f') {
			l.FzfOutput = ""
			launch := tui.LaunchEvent{
				CmdBuilder: func() (cmd *exec.Cmd, captureOutput bool, err error) {
					cmd = exec.Command("fzf")
					home, err := os.UserHomeDir()
					if err != nil {
						return nil, false, err
					}
					cmd.Dir = home
					cmd.Stdin = strings.NewReader("A\nB\nC\nD")

					return cmd, true, nil
				},
				OnResult: func(stdout []byte, err error) tui.Task {
					txt := string(bytes.TrimSpace(stdout))
					l.FzfOutput = fmt.Sprintf("fzf: '%s' err=%v", txt, err)

					return tui.TaskNone()
				},
			}

			return tui.TaskOne(launch)
		}
	}

	return tui.TaskNone()
}

func (l *Launcher) Render(ctx tui.RenderContext) {
	print := func(text string, row int) {
		style := tui.DefaultStyle.
			Fg(tui.ColorHex(0xFFA500)).
			Bg(tui.ColorHex(0x303030))
		for i, r := range []rune(text) {
			ctx.Viewport.SetAt(i, row, r, nil, 1, style)
		}
	}

	text := "n for nvim, f for fzf, ctrl+c/esc for quit"
	print(text, 0)

	if l.NvimOutput != "" {
		print(l.NvimOutput, 1)
	}
	if l.FzfOutput != "" {
		print(l.FzfOutput, 2)
	}
}

func main() {
	app := tui.App{}

	c := &Launcher{}

	if err := app.Run(c); err != nil {
		log.Fatalf("%v", err)
	}
}
