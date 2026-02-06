package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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

	var processList []ProcessInfo
	processListCh := make(chan []ProcessInfo)

	var wg sync.WaitGroup
	wg.Go(func() { eventsF(ctx) })
	wg.Go(func() {
		for {
			list, err := ListProcesses()
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
		layout := view.SplitH(vp, view.Fixed("search", 3), view.Dynamic("body", 3), view.Fixed("help", 1))
		search, body, help := layout["search"], layout["body"], layout["help"]

		draw.Box(search, draw.BoxBorderThin, screen.DefaultStyle.Fg(screen.ColorGreen))
		draw.Line(search.Offset(1), "Type ...", screen.DefaultStyle.Fg(screen.ColorHex(0x0F0F0F)))

		draw.Box(body, draw.BoxBorderDouble, screen.DefaultStyle)

		if processList == nil {
			draw.Line(body.Offset(1), "waiting...", screen.DefaultStyle)
		} else {
			line := fmt.Sprintf("%-10s %-20s %-20s %s\n", "PID", "Kernel Time", "User Time", "Name")
			draw.Line(body.Offset(1, 0, 0, 1), line, screen.DefaultStyle.Fg(screen.ColorGreen))

			vp := body.Offset(2, 0, 0, 1)
			w, h := vp.Size()

			for idx, info := range processList {
				if idx >= h-1 {
					break
				}

				draw.Text(
					vp.Offset(idx, 0, 0, 0).Slice(0, 0, w, 1),
					draw.TextChunk{
						Text:  fmt.Sprintf("%-10d ", info.PID),
						Style: screen.DefaultStyle.Fg(screen.ColorYellow),
					},
					draw.TextChunk{
						Text:  fmt.Sprintf("%-20s ", info.KernelTime.String()),
						Style: screen.DefaultStyle.Fg(screen.ColorBlue),
					},
					draw.TextChunk{
						Text:  fmt.Sprintf("%-20s ", info.UserTime.String()),
						Style: screen.DefaultStyle.Fg(screen.ColorRed),
					},
					draw.TextChunk{
						Text:  info.Name,
						Style: screen.DefaultStyle.Fg(screen.ColorCyan),
					})
			}
		}

		draw.Line(help, "[Q]uit", screen.DefaultStyle)

		buf.WriteTo(terminal)

		select {
		case <-ctx.Done():
			return
		case list := <-processListCh:
			processList = list
		case <-terminal.ResizeC():
			w, h, _ = terminal.GetSize()
			buf.Resize(w, h)
			buf.Flush(terminal)
		case e := <-eventsCh:
			if k, ok := e.(vt.KeyEvent); ok {
				if k.IsKey(vt.KeyCtrlC, vt.KeyQ, vt.KeyEsc) {
					cancel()
				}
			}
		}
	}
}
