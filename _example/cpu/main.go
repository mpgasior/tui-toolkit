package main

import (
	"context"
	"log"
	"strconv"
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
		screen.DefaultStyle = screen.DefaultStyle.Bg(screen.ColorHex(0x3030))
		draw.Clear(vp, screen.DefaultStyle)

		layout := view.SplitH(vp, view.Fixed("search", 3), view.Dynamic("body", 3), view.Fixed("help", 1))
		search, body, help := layout["search"], layout["body"], layout["help"]

		draw.Box(search, draw.BoxBorderThin, screen.DefaultStyle.Fg(screen.ColorGreen))
		draw.Line(search.Offset(1), "Type ...", screen.DefaultStyle.Fg(screen.ColorHex(0x0F0F0F)))

		draw.Box(body, draw.BoxBorderDouble, screen.DefaultStyle)

		if processList == nil {
			draw.Line(body.Offset(1), "waiting...", screen.DefaultStyle)
		} else {
			drawLine := func(vp view.Port, pid draw.TextChunk, name draw.TextChunk, kernel draw.TextChunk, user draw.TextChunk) {
				layout := view.SplitV(vp,
					view.Fixed("pid", 7),
					view.Dynamic("name", 25),
					view.Dynamic("kernel", 5),
					view.Dynamic("user", 5))

				draw.Text(layout["pid"], pid)
				draw.Text(layout["name"], name)
				draw.Text(layout["kernel"], kernel)
				draw.Text(layout["user"], user)
			}

			drawInfo := func(vp view.Port, info ProcessInfo) {
				drawLine(vp,
					draw.TextChunk{
						Text:  strconv.FormatInt(int64(info.PID), 10),
						Style: screen.DefaultStyle.Fg(screen.ColorYellow),
					},
					draw.TextChunk{
						Text:  info.Name,
						Style: screen.DefaultStyle.Fg(screen.ColorCyan),
					},
					draw.TextChunk{
						Text:  info.KernelTime.String(),
						Style: screen.DefaultStyle.Fg(screen.ColorBlue),
					},
					draw.TextChunk{
						Text:  info.UserTime.String(),
						Style: screen.DefaultStyle.Fg(screen.ColorRed),
					})
			}

			headerStyle := screen.DefaultStyle.
				Fg(screen.ColorGreen).
				Attr(screen.AttrUnderline)

			drawLine(body.Offset(1, 0, 0, 1),
				draw.TextChunk{"PID", headerStyle},
				draw.TextChunk{"Name", headerStyle},
				draw.TextChunk{"Kernel", headerStyle},
				draw.TextChunk{"User", headerStyle})

			vp := body.Offset(2, 0, 0, 1)
			w, h := vp.Size()

			for idx, info := range processList {
				if idx >= h-1 {
					break
				}

				row := vp.Offset(idx, 0, 0, 0).Slice(0, 0, w, 1)

				drawInfo(row, info)
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
