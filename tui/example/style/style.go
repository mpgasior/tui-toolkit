package main

import (
	"fmt"

	"github.com/mpgasior/tui-go/tui"
	"github.com/mpgasior/tui-go/vt"
)

func main() {
	styles := []tui.Style{
		tui.DefaultStyle,
		tui.DefaultStyle.
			Fg(tui.ColorBlue).
			Bg(tui.ColorYellow),
		tui.DefaultStyle.
			Attr(tui.AttrBold).
			Bg(tui.ColorRed),
		tui.DefaultStyle.
			Attr(tui.AttrBold).
			Bg(tui.ColorHex(0xFF69B4)),
		tui.DefaultStyle.
			Attr(tui.AttrBold | tui.AttrUnderline).
			Fg(tui.ColorHex(0xFF69B4)),
	}

	for _, style := range styles {
		fmt.Printf("%sTEST%s\r\n", style.Sequence(), vt.SGRReset)
	}
}
