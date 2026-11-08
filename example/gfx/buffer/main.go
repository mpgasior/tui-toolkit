package main

import (
	"context"
	"log"
	"os"

	"github.com/mpgasior/tui-toolkit/gfx"
	"github.com/mpgasior/tui-toolkit/termio"
)

func main() {
	mods := []termio.SessionModifier{
		termio.WithAltScreen,
	}

	term, err := termio.OpenSession(os.Stdin, os.Stdout, mods...)
	if err != nil {
		log.Fatalln(err)
	}
	defer term.Close()

	width, height, err := term.GetSize()
	if err != nil {
		log.Fatalln(err)
	}

	buf := gfx.NewBuffer(width, height)
	lineH := func(y int, r rune, s gfx.Style) {
		for x := 0; x < buf.Width; x += 1 {
			buf.Set(x, y, r, s)
		}
	}

	lineH(0, '-', gfx.StyleDefault.SetFgColor(gfx.Red))
	lineH(height-1, '-', gfx.StyleDefault.SetFgColor(gfx.Green))

	if _, err := gfx.WriteBufferDiff(term, nil, buf); err != nil {
		log.Fatalln(err)
	}
	gfx.WriteReset(term)

	var p [1]byte
	_, err = term.ReadContext(context.Background(), p[:])
	if err != nil {
		log.Fatalln(err)
	}
}
