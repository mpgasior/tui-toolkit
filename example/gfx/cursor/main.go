package main

import (
	"context"
	"log"
	"os"

	"github.com/mpgasior/tui-toolkit/ansi"
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

	shape := []gfx.CursorShape{
		gfx.CursorShapeBlock,
		gfx.CursorShapeUnderline,
		gfx.CursorShapeBar,
	}

	blink := []bool{true, false}

	ctx := context.Background()
	var p [1]byte
	for _, s := range shape {
		for _, b := range blink {
			cs := gfx.CursorStyleDefault.SetShape(s).SetBlink(b).SetColor(gfx.RGB(0xFF000F))
			buf.SetCursor(0, 0, cs)

			if _, err := gfx.WriteBufferDiff(term, nil, buf); err != nil {
				return err
			}
			if _, err := term.ReadContext(ctx, p[:]); err != nil {
				return err
			}
		}
	}

	dst := ansi.AppendResetCursorColorRGB(nil)
	dst = ansi.AppendCursorStyle(dst, ansi.CursorDefault)
	if _, err := term.Write(dst); err != nil {
		return err
	}

	return nil
}
