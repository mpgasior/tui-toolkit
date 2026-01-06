package vt

import (
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

// Text Formatting
// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#text-formatting
const (
	SGRFmt   = CSI + "%dm"
	SGRReset = CSI + "0m"
)

func FormatSGR[T constraints.Integer](p ...T) string {
	if len(p) == 0 {
		return SGRReset
	}

	parts := make([]string, len(p))
	for i, part := range p {
		parts[i] = strconv.FormatInt(int64(part), 10)
	}

	return CSI + strings.Join(parts, ";") + "m"
}

type Attr uint32

// SGR Parameter Codes
const (
	AttrReset Attr = iota
	AttrBold
	AttrFaint
	AttrItalic
	AttrUnderline
	AttrBlinkSlow
	AttrBlinkRapid
	AttrReverseVideo
	AttrConceal
	AttrStrikethrough
)

type Color int

const (
	ColorBlack Color = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	_
	ColorReset
)

type FgColor uint32

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
	_
	FgDefault
)

type FgBrightColor uint32

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

type BgColor uint32

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
	_
	BgDefault
)

type BgBrightColor uint32

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
