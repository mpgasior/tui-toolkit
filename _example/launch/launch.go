package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/session"
	"github.com/mpgasior/tui-toolkit/termx"
	"github.com/mpgasior/tui-toolkit/view"
	"github.com/mpgasior/tui-toolkit/vt"
)

func main() {
	tty, err := termx.OpenTTY()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	terminal, err := termx.New(tty.In, tty.Out)
	if err != nil {
		log.Fatal(err)
	}
	defer terminal.Close()

	eventsF, eventsCh := vt.Events(terminal)

	ses := session.New(
		terminal.MakeRaw,
		func() (restore func() error, err error) {
			return vt.EnterMode(terminal, vt.ModeAlternateScreen)
		},
		func() (restore func() error, err error) {
			return vt.EnterMode(terminal, vt.ModeBracketedPaste)
		},
		func() (restore func() error, err error) {
			return vt.ExitMode(terminal, vt.ModeShowCursor)
		},
		func() (restore func() error, err error) {
			ctx, cancel := context.WithCancel(context.Background())
			var wg sync.WaitGroup
			wg.Go(func() { eventsF(ctx) })

			return func() error {
				cancel()
				wg.Wait()
				return nil
			}, nil
		},
	)

	if err := ses.Start(); err != nil {
		log.Fatal(err)
	}
	defer ses.Stop()

	w, h, err := terminal.GetSize()
	if err != nil {
		log.Fatal(err)
	}

	buf := screen.New(w, h)

	for {
		vp := view.NewPort(buf)

		draw.Fill(vp, ' ', screen.DefaultStyle)
		draw.Box(vp, draw.BoxBorderThin, screen.DefaultStyle.Fg(screen.ColorGreen))
		draw.Box(vp.Offset(10), draw.BoxBorderDouble, screen.DefaultStyle.Fg(screen.ColorBlue))
		w, h := vp.Size()
		draw.Text(vp.Offset(11), fmt.Sprintf("%d-%d", w, h), screen.DefaultStyle)

		buf.WriteTo(terminal)

		select {
		case <-terminal.ResizeC():
			w, h, _ = terminal.GetSize()
			buf.Resize(w, h)
			buf.Flush(terminal)
		case e := <-eventsCh:
			if k, ok := e.(vt.KeyEvent); ok {
				if k.IsKey(vt.KeyCtrlC, vt.KeyQ, vt.KeyEsc) {
					return
				}

				if k.IsKey(vt.KeyE, vt.KeyO) {
					ses.RunSuspended(func() error {
						cmd := exec.Command("nvim")
						cmd.Stdin = tty.In
						cmd.Stdout = tty.Out
						cmd.Stderr = os.Stderr

						return cmd.Run()
					})
				}
			}
		}
	}
}
