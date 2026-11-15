package main

import (
	"context"
	"log"
	"os"

	"github.com/mpgasior/tui-toolkit/gfx"
	"github.com/mpgasior/tui-toolkit/termio"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	mods := []termio.SessionModifier{
		termio.WithAltScreen,
	}
	term, err := termio.OpenSession(os.Stdin, os.Stdout, mods...)
	if err != nil {
		return err
	}
	defer term.Close()

	w, h, err := term.GetSize()
	if err != nil {
		return err
	}

	screen := gfx.NewScreen(w, h, nil)

	styles := []gfx.Style{
		gfx.StyleDefault.SetBgColor(gfx.Red),
		gfx.StyleDefault.SetBgColor(gfx.Green),
		gfx.StyleDefault.SetBgColor(gfx.Blue),
	}

	idx := 0
	nextStyle := func() gfx.Style {
		idx += 1
		return styles[idx%len(styles)]
	}

	for {
		var p [64]byte
		_, err := term.ReadContext(context.Background(), p[:])
		if err != nil {
			return err
		}

		if p == [64]byte{3} {
			return nil
		}

		if p == [64]byte{' '} {
			Fill(screen.Back, nextStyle())
		}

		screen.Swap()
		if _, err := screen.WriteTo(term); err != nil {
			return err
		}
	}
}

func Fill(buf *gfx.Buffer, s gfx.Style) {
	for y := 0; y < buf.Height; y += 1 {
		rowOffset := y * buf.Width
		for x := 0; x < buf.Width; x += 1 {
			pos := rowOffset + x
			buf.Cells[pos] = gfx.Cell{R: ' ', Style: s}
		}
	}
}
