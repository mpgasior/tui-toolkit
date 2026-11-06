package main

import (
	"fmt"
	"os"

	"github.com/mpgasior/tui-toolkit/ansi"
)

func main() {
	var buf []byte
	buf = ansi.AppendSGRStart(buf)
	buf = ansi.AppendSGRParam(buf, ansi.Bold)
	buf = ansi.AppendSGRParam(buf, ansi.Underline)
	buf = ansi.AppendSGRParam(buf, ansi.FgBlue)
	buf = ansi.AppendSGRParam(buf, ansi.BgBrightCyan)
	buf = ansi.AppendSGREnd(buf)

	os.Stdout.Write(buf)
	os.Stdout.WriteString("A bold and underline blue text on bright cyan background.")
	buf = ansi.AppendSGRReset(nil)
	os.Stdout.Write(buf)
	fmt.Fprintln(os.Stdout)

	buf = ansi.AppendSGRStart(nil)
	buf = ansi.AppendSGRRGBFg(buf, 0xFF0F00)
	buf = ansi.AppendSGRRGBBg(buf, 0xCCCCCC)
	buf = ansi.AppendSGREnd(buf)

	os.Stdout.Write(buf)
	fmt.Fprint(os.Stdout, "An orange text on a gray background.")
	buf = ansi.AppendSGRReset(nil)
	os.Stdout.Write(buf)
	fmt.Fprintln(os.Stdout)
}
