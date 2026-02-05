package main

import (
	"log"

	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/termx"
	"github.com/mpgasior/tui-toolkit/view"
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

	w, h, err := terminal.GetSize()
	if err != nil {
		log.Fatal(err)
	}

	buf := screen.New(w, h)
	vp := view.NewPort(buf)

	draw.Box(vp, draw.BoxBorderThin, screen.DefaultStyle.Fg(screen.ColorGreen))
	draw.Box(vp.Offset(10), draw.BoxBorderDouble, screen.DefaultStyle.Fg(screen.ColorBlue))

	buf.WriteTo(terminal)
}
