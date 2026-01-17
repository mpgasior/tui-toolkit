package vt

// Text Formatting
// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#text-formatting
const (
	SGRFmt = CSI + "%dm"
)

// SGR Parameter Codes
const (
	AttrReset         = 0
	AttrBold          = 1
	AttrFaint         = 2
	AttrItalic        = 3
	AttrUnderline     = 4
	AttrBlinkSlow     = 5
	AttrBlinkRapid    = 6
	AttrReverseVideo  = 7
	AttrConceal       = 8
	AttrStrikethrough = 9
)

// Foreground Colors (30-37)
const (
	FgBlack = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
	FgDefault = 39
)

// Foreground Colors (Bright/High Intensity)
const (
	FgBrightBlack = iota + 90
	FgBrightRed
	FgBrightGreen
	FgBrightYellow
	FgBrightBlue
	FgBrightMagenta
	FgBrightCyan
	FgBrightWhite
)

// Background Colors (40-47)
const (
	BgBlack = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
	BgDefault = 49
)

// Background Colors (Bright/High Intensity)
const (
	BgBrightBlack = iota + 100
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
