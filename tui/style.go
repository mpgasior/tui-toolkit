package tui

import (
	"github.com/mpgasior/tui-go/vt"
)

type Attr uint8

const (
	AttrReset Attr = 1 << iota
	AttrBold
	AttrFaint
	AttrItalic
	AttrUnderline
	AttrBlinkSlow
	AttrBlinkRapid
	AttrStrikethrough
)

type Color uint32

const (
	ColorModeANSI    uint32 = 0 << 24
	ColorModePalette uint32 = 1 << 24
	ColorModeRGB     uint32 = 2 << 24
)

func ColorVt(c vt.Color) Color {
	return Color(uint32(c) | ColorModeANSI)
}

func ColorPalette(v uint8) Color {
	return Color(uint32(v) | ColorModePalette)
}

func ColorRGB(r, g, b uint8) Color {
	return Color(uint32(r)<<16 | uint32(g)<<8 | uint32(b) | ColorModeRGB)
}

func ColorHex(hex uint32) Color {
	return Color(hex&0xFFFFFF | ColorModeRGB)
}

var (
	ColorDefault = ColorReset
	ColorReset   = ColorVt(vt.ColorReset)

	ColorBlack   = ColorVt(vt.ColorBlack)
	ColorRed     = ColorVt(vt.ColorRed)
	ColorGreen   = ColorVt(vt.ColorGreen)
	ColorYellow  = ColorVt(vt.ColorYellow)
	ColorBlue    = ColorVt(vt.ColorBlue)
	ColorMagenta = ColorVt(vt.ColorMagenta)
	ColorCyan    = ColorVt(vt.ColorCyan)
	ColorWhite   = ColorVt(vt.ColorWhite)

	ColorBrightBlack   = ColorVt(vt.ColorBrightBlack)
	ColorBrightRed     = ColorVt(vt.ColorBrightRed)
	ColorBrightGreen   = ColorVt(vt.ColorBrightGreen)
	ColorBrightYellow  = ColorVt(vt.ColorBrightYellow)
	ColorBrightBlue    = ColorVt(vt.ColorBrightBlue)
	ColorBrightMagenta = ColorVt(vt.ColorBrightMagenta)
	ColorBrightCyan    = ColorVt(vt.ColorBrightCyan)
	ColorBrightWhite   = ColorVt(vt.ColorBrightWhite)
)

const (
	shiftAttr uint8 = 0
	shiftFg   uint8 = 8
	shiftBg         = 36

	bitsPerColor        = 28
	maskAttr     uint64 = 0xFF
	maskColor    uint64 = 0xFFFFFFF
)

type Style uint64

var (
	DefaultStyle = NewStyle()
)

func NewStyle() Style {
	s := Style(0).
		Fg(ColorReset).
		Bg(ColorReset)

	return s
}

func (s Style) Fg(c Color) Style {
	cleanMask := ^(maskColor << shiftFg)
	return (s & Style(cleanMask)) | Style(uint64(c))<<shiftFg
}

func (s Style) Bg(c Color) Style {
	cleanMask := ^(maskColor << shiftBg)
	return (s & Style(cleanMask)) | Style(uint64(c))<<shiftBg
}

func (s Style) Attr(a Attr) Style {
	return s | Style(a)
}

func (s Style) NoAttr(a Attr) Style {
	return s & ^(Style(a) & Style(maskAttr))
}

var attrMap = map[Attr]vt.Attr{
	AttrReset:         vt.AttrReset,
	AttrBold:          vt.AttrBold,
	AttrFaint:         vt.AttrFaint,
	AttrItalic:        vt.AttrItalic,
	AttrUnderline:     vt.AttrUnderline,
	AttrBlinkSlow:     vt.AttrBlinkSlow,
	AttrBlinkRapid:    vt.AttrBlinkRapid,
	AttrStrikethrough: vt.AttrStrikethrough,
}
