package gfx

type Attr byte

const (
	AttrNone      Attr = 0
	AttrBold      Attr = 1 << 0
	AttrDim       Attr = 1 << 1
	AttrItalic    Attr = 1 << 2
	AttrUnderline Attr = 1 << 3
)

var (
	Bold      = AttrBold
	Dim       = AttrDim
	Italic    = AttrItalic
	Underline = AttrUnderline
)
