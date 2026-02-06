package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/_example/cpu/ui"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
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

	restoreX, err := terminal.MakeRaw()
	if err != nil {
		log.Fatal(err)
	}
	defer restoreX()

	eventsF, eventsCh := vt.Events(terminal)

	w, h, err := terminal.GetSize()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var procTable ui.ProcessTable
	var searchInput ui.SearchInput

	processListCh := make(chan []process.ProcessInfo)

	var wg sync.WaitGroup
	wg.Go(func() { eventsF(ctx) })
	wg.Go(func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
		for {
			list, err := process.List()
			if err != nil {
				break
			}

			select {
			case <-ctx.Done():
				return
			case processListCh <- list:
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
	})

	buf := screen.New(w, h)
	restore, err := vt.EnterMode(terminal, vt.ModeAlternateScreen)
	if err != nil {
		log.Fatal(err)
	}
	defer restore()

	for {
		vp := view.NewPort(buf)
		screen.DefaultStyle = screen.DefaultStyle.
			Bg(screen.ColorHex(0x2B2D42)).
			Fg(screen.ColorHex(0xEDF2F4)).
			Fg(screen.ColorHex(0xEF233C)).
			Fg(screen.ColorHex(0x8D99AE))

		draw.Clear(vp, screen.DefaultStyle)

		layout := view.SplitH(vp, view.Fixed("search", 3), view.Dynamic("body", 3), view.Fixed("help", 1))
		search, body, help := layout["search"], layout["body"], layout["help"]

		searchInput.Draw(search)
		procTable.Draw(body)

		draw.Line(help, "[Q]uit", screen.DefaultStyle)

		buf.WriteTo(terminal)

		select {
		case <-ctx.Done():
			return
		case snapshot := <-processListCh:
			procTable.Rows = snapshot
		case <-terminal.ResizeC():
			w, h, _ = terminal.GetSize()
			buf.Resize(w, h)
			buf.Flush(terminal)
		case e := <-eventsCh:
			if k, ok := e.(vt.KeyEvent); ok {
				if k.IsKey(vt.KeyCtrlC, vt.KeyQ, vt.KeyEsc) {
					cancel()
				}
				searchInput.Update(k)
			}
		}
	}
}
