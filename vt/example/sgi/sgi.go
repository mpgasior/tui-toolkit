package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/mpgasior/tui-go/vt"
	"golang.org/x/exp/constraints"
)

func main() {
	bgStandard := []vt.BgColor{
		vt.BgBlack, vt.BgRed,
		vt.BgGreen, vt.BgYellow,
		vt.BgBlue, vt.BgMagenta,
		vt.BgCyan, vt.BgWhite,
	}

	printRow("BG Std", bgStandard)

	bgBright := []vt.BgBrightColor{
		vt.BgBrightBlack, vt.BgBrightRed,
		vt.BgBrightGreen, vt.BgBrightYellow,
		vt.BgBrightBlue, vt.BgBrightMagenta,
		vt.BgBrightCyan, vt.BgBrightWhite,
	}

	printRow("BG Brgt", bgBright)

	fmt.Println()

	fgStandard := []vt.FgColor{
		vt.FgBlack, vt.FgRed,
		vt.FgGreen, vt.FgYellow,
		vt.FgBlue, vt.FgMagenta,
		vt.FgCyan, vt.FgWhite,
	}

	printRow("FG Std", fgStandard)

	fgBright := []vt.FgBrightColor{
		vt.FgBrightBlack, vt.FgBrightRed,
		vt.FgBrightGreen, vt.FgBrightYellow,
		vt.FgBrightBlue, vt.FgBrightMagenta,
		vt.FgBrightCyan, vt.FgBrightWhite,
	}

	printRow("FG Brgt", fgBright)

	attrs := []vt.Attr{
		vt.AttrReset, vt.AttrBold,
		vt.AttrFaint, vt.AttrItalic,
		vt.AttrUnderline, vt.AttrBlinkSlow,
		vt.AttrBlinkRapid, vt.AttrReverseVideo,
		vt.AttrConceal, vt.AttrStrikethrough,
	}

	printRow("Attrs", attrs)

	fmt.Println()
	reset := fmt.Sprintf(vt.SGRFmt, vt.AttrReset)
	for i := range 256 {
		fmt256 := fmt.Sprintf(vt.BgColor256Fmt, i)
		sgi := fmt.Sprint(vt.CSI, fmt256, "m")
		fmt.Printf("%-3d %s %s ", i, sgi, reset)

		if (i+1)%16 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func printRow[T constraints.Integer](label string, items []T) {
	reset := fmt.Sprintf(vt.SGRFmt, vt.AttrReset)
	fmt.Printf("%-8s | ", label)

	for _, item := range items {
		sgi := fmt.Sprintf(vt.SGRFmt, item)
		str := strconv.FormatInt(int64(item), 10)
		text := lastCapitalized(str)
		fmt.Printf("%s %s %s ", sgi, text, reset)
	}
	fmt.Println()
}

func lastCapitalized(s string) string {
	i := strings.LastIndexFunc(s, unicode.IsUpper)
	if i == -1 {
		return s
	}

	return s[i:]
}
