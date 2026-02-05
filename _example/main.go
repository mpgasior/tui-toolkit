package main

import (
	"context"
	"fmt"
	"log"
	"sync"

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

	var wg sync.WaitGroup
	wg.Go(func() { eventsF(ctx) })

	buf := screen.New(w, h)
	restore, err := vt.EnterMode(terminal, vt.ModeAlternateScreen)
	if err != nil {
		log.Fatal(err)
	}
	defer restore()

	for {
		vp := view.NewPort(buf)

		draw.Fill(vp, ' ', screen.DefaultStyle)
		draw.Box(vp, draw.BoxBorderThin, screen.DefaultStyle.Fg(screen.ColorGreen))
		draw.Box(vp.Offset(10), draw.BoxBorderDouble, screen.DefaultStyle.Fg(screen.ColorBlue))
		w, h := vp.Size()
		draw.Text(vp.Offset(11), fmt.Sprintf("%d-%d", w, h), screen.DefaultStyle)

		buf.WriteTo(terminal)

		select {
		case <-ctx.Done():
			return
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
