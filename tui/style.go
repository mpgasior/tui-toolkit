package tui

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
	ColorModeDefault uint32 = 3 << 24
)

const (
	shiftAttr uint8 = 0
	shiftFg   uint8 = 8
	shiftBg         = 36

	bitsPerColor        = 28
	maskAttr     uint64 = 0xFF
	maskColor    uint64 = 0xFFFFFFF
)

func ANSI(v uint8) Color {
	return Color(uint32(v) | ColorModeANSI)
}

func Palette(v uint8) Color {
	return Color(uint32(v) | ColorModePalette)
}

func RGB(r, g, b uint8) Color {
	return Color(uint32(r)<<16 | uint32(g)<<8 | uint32(b) | ColorModeRGB)
}

type Style uint64

func EmptyStyle() Style {
	var s uint64

	s |= uint64(ColorModeDefault) << shiftFg
	s |= uint64(ColorModeDefault) << shiftBg

	return Style(s)
}

func (s Style) FgColor(c Color) Style {
	cleanMask := ^(maskColor << shiftFg)
	return (s & Style(cleanMask)) | Style(uint64(c))<<shiftFg
}

func (s Style) BgColor(c Color) Style {
	cleanMask := ^(maskColor << shiftBg)
	return (s & Style(cleanMask)) | Style(uint64(c))<<shiftBg
}
