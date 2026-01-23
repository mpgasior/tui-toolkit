package vt

// Text Formatting
// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#text-formatting
const (
	SGRFmt = CSI + "%dm"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=Attr
type Attr int

// SGR Parameter Codes
const (
	AttrReset         Attr = 0
	AttrBold          Attr = 1
	AttrFaint         Attr = 2
	AttrItalic        Attr = 3
	AttrUnderline     Attr = 4
	AttrBlinkSlow     Attr = 5
	AttrBlinkRapid    Attr = 6
	AttrReverseVideo  Attr = 7
	AttrConceal       Attr = 8
	AttrStrikethrough Attr = 9
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=FgColor
type FgColor int

// Foreground Colors (30-37)
const (
	FgBlack FgColor = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
	FgDefault = 39
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=FgBrightColor
type FgBrightColor int

// Foreground Colors (Bright/High Intensity)
const (
	FgBrightBlack FgBrightColor = iota + 90
	FgBrightRed
	FgBrightGreen
	FgBrightYellow
	FgBrightBlue
	FgBrightMagenta
	FgBrightCyan
	FgBrightWhite
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=BgColor
type BgColor int

// Background Colors (40-47)
const (
	BgBlack BgColor = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
	BgDefault = 49
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=BgBrightColor
type BgBrightColor int

// Background Colors (Bright/High Intensity)
const (
	BgBrightBlack BgBrightColor = iota + 100
	BgBrightRed
	BgBrightGreen
	BgBrightYellow
	BgBrightBlue
	BgBrightMagenta
	BgBrightCyan
	BgBrightWhite
)

// Extended Colors
const (
	FgColor256Fmt = "38;5;%d"
	BgColor256Fmt = "48;5;%d"

	FgColorRGBFmt = "38;2;%d;%d;%d"
	BgColorRGBFmt = "48;2;%d;%d;%d"
)

// Screen Colors
const (
	ModifyScreenColorsFmt = OSC + "4;%d;rgb:%02x/%02x/%02x" + ST
)
