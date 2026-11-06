package gfx

type Style uint64

const (
	fgColorShift uint64 = 0
	bgColorShift uint64 = 26
	attrShift    uint64 = 26 + 26

	colorMask uint64 = 0x3FFFFFF
	attrMask  uint64 = 0xF

	fgMask uint64 = colorMask << fgColorShift
	bgMask uint64 = colorMask << bgColorShift
)

func (s Style) SetFgColor(c Color) Style {
	style := uint64(s) & ^fgMask
	style |= (uint64(c) << fgColorShift)

	return Style(style)
}

func (s Style) SetBgColor(c Color) Style {
	style := uint64(s) & ^bgMask
	style |= (uint64(c) << bgColorShift)

	return Style(style)
}

func (s Style) SetAttr(a Attr) Style {
	style := uint64(s) &^ (attrMask << attrShift)
	style |= ((uint64(a) & attrMask) << attrShift)

	return Style(style)
}

func (s Style) Unpack() (fg, bg Color, attr Attr) {
	fg = Color((uint64(s) >> fgColorShift) & colorMask)
	bg = Color((uint64(s) >> bgColorShift) & colorMask)
	attr = Attr((uint64(s) >> attrShift) & attrMask)

	return fg, bg, attr
}

var (
	StyleDefault = Style(0).
		SetBgColor(ColorDefault).
		SetFgColor(ColorDefault).
		SetAttr(AttrNone)
)
