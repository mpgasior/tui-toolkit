package screen

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

	MaskColorMode  uint32 = 0xFF << 24
	MaskColorValue uint32 = 0x00FFFFFF
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
)

const (
	bitsPerColor       = 28
	shiftAttr    uint8 = 0
	ShiftFg      uint8 = 8
	ShiftBg            = ShiftFg + bitsPerColor

	MaskAttr    uint64 = 0xFF
	MaskColor   uint64 = 0xFFFFFFF
	MaskFgColor uint64 = MaskColor << ShiftFg
	MaskBgColor uint64 = MaskColor << ShiftBg
)

type Style uint64

var (
	DefaultStyle = Style(0).
		Fg(ColorReset).
		Bg(ColorReset)
)

func (s Style) Fg(c Color) Style {
	cleanMask := ^(MaskColor << ShiftFg)
	return (s & Style(cleanMask)) | Style(uint64(c))<<ShiftFg
}

func (s Style) Bg(c Color) Style {
	cleanMask := ^(MaskColor << ShiftBg)
	return (s & Style(cleanMask)) | Style(uint64(c))<<ShiftBg
}

func (s Style) Attr(a Attr) Style {
	return s | Style(a)
}

func (s Style) NoAttr(a Attr) Style {
	return s & ^(Style(a) & Style(MaskAttr))
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

func (a Attr) Params() []uint32 {
	var parts []uint32
	for k, v := range attrMap {
		if a&k != 0 {
			parts = append(parts, uint32(v))
		}
	}
	return parts
}

func (c Color) Params(isBg bool) []uint32 {
	mode := uint32(c) & MaskColorMode
	val := uint32(c) & MaskColorValue

	var bgOffset uint32
	if isBg {
		bgOffset += 10
	}

	switch mode {
	case ColorModeANSI:
		return []uint32{val + 30 + bgOffset}
	case ColorModePalette:
		return []uint32{38 + bgOffset, 5, val}
	case ColorModeRGB:
		return []uint32{
			38 + bgOffset, 2,
			(val >> 16) & 0xFF,
			(val >> 8) & 0xFF,
			(val) & 0xFF,
		}
	default:
		return nil
	}
}

func (s Style) Sequence() string {
	attr := Attr(uint64(s) & MaskAttr)
	fg := Color((uint64(s) & MaskFgColor) >> ShiftFg)
	bg := Color((uint64(s) & MaskBgColor) >> ShiftBg)

	var params []uint32
	params = append(params, attr.Params()...)
	params = append(params, fg.Params(false)...)
	params = append(params, bg.Params(true)...)

	return vt.FormatSGR(params...)
}
