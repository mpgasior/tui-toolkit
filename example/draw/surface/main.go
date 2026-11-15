package main

import (
	"context"
	"log"
	"os"

	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/gfx"
	"github.com/mpgasior/tui-toolkit/termio"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	term, err := termio.OpenSession(os.Stdin, os.Stdout, termio.WithAltScreen)
	if err != nil {
		return err
	}
	defer term.Close()

	w, h, err := term.GetSize()
	if err != nil {
		return err
	}

	buf := gfx.NewBuffer(w, h)
	sf := draw.NewSurface(buf)

	colors := []gfx.Color{
		gfx.Black,
		gfx.Red,
		gfx.Green,
		gfx.Yellow,
		gfx.Blue,
		gfx.Magenta,
		gfx.Cyan,
		gfx.White,
	}

	left := sf.Sub(0, 0, w/2, h)
	right := sf.Sub(w/2, 0, w/2, h)

	fill := func(sf *draw.Surface, r rune, s gfx.Style) {
		for row := range sf.Height {
			for col := range sf.Width {
				sf.Set(col, row, r, s)
			}
		}
	}

	for i := range min(sf.Width, sf.Height) {
		color := colors[i%len(colors)]
		style := gfx.StyleDefault.SetBgColor(color)
		sf := left.Inset(i, i, i, i)
		fill(sf, ' ', style)
	}

	fill(right, ' ', gfx.StyleDefault.SetBgColor(gfx.Green))
	if _, err = gfx.WriteBufferDiff(term, nil, buf); err != nil {
		return err
	}

	ctx := context.Background()
	var p [1]byte
	_, err = term.ReadContext(ctx, p[:])

	return err
}
