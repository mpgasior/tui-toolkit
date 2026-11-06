package gfx

import "github.com/mpgasior/tui-toolkit/ansi"

func (a Attr) AppendTo(dst []byte) []byte {
	if a == AttrNone {
		return dst
	}

	if a&AttrBold != 0 {
		dst = ansi.AppendSGRParam(dst, ansi.Bold)
	}
	if a&AttrDim != 0 {
		dst = ansi.AppendSGRParam(dst, ansi.Dim)
	}
	if a&AttrItalic != 0 {
		dst = ansi.AppendSGRParam(dst, ansi.Italic)
	}
	if a&AttrUnderline != 0 {
		dst = ansi.AppendSGRParam(dst, ansi.Underline)
	}
	return dst
}

func (c Color) AppendFgTo(dst []byte) []byte {
	switch c.Kind() {
	case ColorKindANSI:
		slot, bright := c.ANSI()
		base := ansi.FgBlack
		if bright {
			base = ansi.FgBrightBlack
		}

		return ansi.AppendSGRParam(dst, int(base)+int(slot))
	case ColorKindRGB:
		return ansi.AppendSGRRGBFg(dst, c.RGB())
	case ColorKindDefault:
		fallthrough
	default:
		return ansi.AppendSGRParam(dst, ansi.FgDefault)
	}
}

func (c Color) AppendBgTo(dst []byte) []byte {
	switch c.Kind() {
	case ColorKindANSI:
		slot, bright := c.ANSI()
		base := ansi.BgBlack
		if bright {
			base = ansi.BgBrightBlack
		}

		return ansi.AppendSGRParam(dst, int(base)+int(slot))
	case ColorKindRGB:
		return ansi.AppendSGRRGBBg(dst, c.RGB())
	case ColorKindDefault:
		fallthrough
	default:
		return ansi.AppendSGRParam(dst, ansi.BgDefault)
	}
}
