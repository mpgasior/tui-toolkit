package main

import (
	"context"
	"fmt"
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

	draw.Clear(sf)
	draw.Fill(sf, ' ', gfx.StyleDefault.SetBgColor(gfx.RGB(0x0F0F0F)))
	draw.Box(sf, draw.BoxOptions{
		Border: draw.BoxBorderThin,
		Style:  gfx.StyleDefault.SetFgColor(gfx.BrightGreen),
	})
	draw.String(sf.Inset(1, 1, 1, 1), "abdcdasdad\r00000\nsadsadas")

	cyan := gfx.StyleDefault.SetFgColor(gfx.Cyan)
	writer := draw.NewTextWriter(sf, draw.WithTextWriterStyle(cyan))
	fmt.Fprintf(writer, "%s -> %s", "hello", "world!")

	if _, err := gfx.WriteBufferDiff(term, nil, buf); err != nil {
		return err
	}

	ctx := context.Background()
	var p [1]byte
	_, err = term.ReadContext(ctx, p[:])

	return err
}
