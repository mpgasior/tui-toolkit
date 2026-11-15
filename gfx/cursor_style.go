package gfx

type CursorShape byte

const (
	CursorShapeBlock     CursorShape = 1
	CursorShapeUnderline CursorShape = 3
	CursorShapeBar       CursorShape = 5
)

const (
	curShapeMask  uint32 = 0xF
	curShapeShift uint32 = 26
	curBlinkFlag  uint32 = 1 << 31
)

type CursorStyle uint32

func (cs CursorStyle) SetColor(c Color) CursorStyle {
	style := uint32(cs) & ^uint32(colorMask)
	style |= uint32(c)

	return CursorStyle(style)
}

func (cs CursorStyle) SetShape(s CursorShape) CursorStyle {
	style := uint32(cs) & ^curShapeMask
	style |= uint32(s) << curShapeShift
	return CursorStyle(style)
}

func (cs CursorStyle) SetBlink(on bool) CursorStyle {
	style := uint32(cs) & ^curBlinkFlag
	if on {
		style |= curBlinkFlag
	}

	return CursorStyle(style)
}

func (cs CursorStyle) Unpack() (color Color, shape CursorShape, blink bool) {
	c := uint32(cs) & uint32(colorMask)
	s := (uint32(cs) >> curShapeShift) & curShapeMask
	blink = uint32(cs)&curBlinkFlag != 0

	return Color(c), CursorShape(s), blink
}

const CursorStyleDefault = CursorStyle(0)
