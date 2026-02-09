package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/mpgasior/tui-toolkit/_example/cpu/process"
	"github.com/mpgasior/tui-toolkit/_example/cpu/ui"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/termx"
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

	u := ui.New()

	processListCh := make(chan []process.ProcessInfo)
	filterCh := make(chan []process.ProcessInfo)

	var wg sync.WaitGroup
	wg.Go(func() { eventsF(ctx) })
	wg.Go(func() {
		process.RefreshLoop(ctx, processListCh, time.Second)
	})
	wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				return
			case snapshot := <-processListCh:
				var filtered []process.ProcessInfo

				for _, p := range snapshot {
					if strings.Contains(p.Name, u.Term()) {
						filtered = append(filtered, p)
					}
				}

				select {
				case <-ctx.Done():
					return
				case filterCh <- filtered:
				}
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
		screen.DefaultStyle = screen.DefaultStyle.
			Bg(screen.ColorHex(0x2B2D42)).
			Fg(screen.ColorHex(0xEDF2F4)).
			Fg(screen.ColorHex(0xEF233C)).
			Fg(screen.ColorHex(0x8D99AE))

		u.Draw(buf)

		buf.WriteTo(terminal)

		select {
		case <-ctx.Done():
			return
		case snapshot := <-filterCh:
			u.OnNewRows(snapshot)
		case <-terminal.ResizeC():
			w, h, _ = terminal.GetSize()
			buf.Resize(w, h)
			buf.Flush(terminal)
		case e := <-eventsCh:
			if k, ok := e.(vt.KeyEvent); ok {
				if k.IsKey(vt.KeyCtrlC) {
					cancel()
				}
				u.HandleKey(k)
			}
		}
	}
}
