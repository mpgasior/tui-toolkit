package main

import (
	"fmt"
	"os"

	"github.com/mpgasior/tui-toolkit/gfx"
)

func main() {
	defer gfx.WriteReset(os.Stdout)
	style := gfx.StyleDefault.
		SetFgColor(gfx.Blue).
		SetBgColor(gfx.BrightCyan).
		SetAttr(gfx.Bold | gfx.Underline)

	gfx.WriteStyleDiff(os.Stdout, gfx.StyleDefault, style)
	fmt.Fprint(os.Stdout, "A bold and underline blue text on bright cyan background.")
	gfx.WriteReset(os.Stdout)
	fmt.Fprintln(os.Stdout)

	style = style.SetFgColor(gfx.RGB(0xFF0F00)).SetBgColor(gfx.RGB(0xCCCCCC)).SetAttr(gfx.AttrNone)
	gfx.WriteStyleDiff(os.Stdout, gfx.StyleDefault, style)
	fmt.Fprintf(os.Stdout, "An orange text on the gray background")
	gfx.WriteReset(os.Stdout)
	fmt.Fprintln(os.Stdout)
}
