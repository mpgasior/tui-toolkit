package gfx

type ColorSlot uint8

const (
	ColorSlotBlack ColorSlot = iota
	ColorSlotRed
	ColorSlotGreen
	ColorSlotYellow
	ColorSlotBlue
	ColorSlotMagenta
	ColorSlotCyan
	ColorSlotWhite
)

type ColorKind uint8

const (
	ColorKindDefault ColorKind = iota
	ColorKindANSI
	ColorKindRGB
)

type Color uint32

const (
	ColorDefault Color = 0

	colorKindMask     uint32 = 0x3
	colorRGBMask      uint32 = 0xFFFFFF
	colorANSISlotMask uint32 = 0x7
	colorKindShift           = 24
	colorBrightShift         = 31

	colorBrightPayloadBit = 1 << 23
)

func ColorRGB(hex uint32) Color {
	return Color(hex&colorRGBMask | uint32(ColorKindRGB)<<colorKindShift)
}

func ColorANSI(slot ColorSlot, bright bool) Color {
	p := uint32(slot)
	p |= uint32(ColorKindANSI) << colorKindShift
	if bright {
		p |= colorBrightPayloadBit
	}
	return Color(p)
}

func (c Color) Kind() ColorKind {
	return ColorKind(uint32(c>>colorKindShift) & colorKindMask)
}

func (c Color) Payload() uint32 {
	return uint32(c) & colorRGBMask
}

func (c Color) ANSI() (slot ColorSlot, bright bool) {
	p := c.Payload()
	slot = ColorSlot(p & colorANSISlotMask)
	bright = (uint32(c) & colorBrightPayloadBit) != 0
	return slot, bright
}

func (c Color) RGB() (hex uint32) {
	return c.Payload()
}

var (
	RGB = ColorRGB

	Black   = ColorANSI(ColorSlotBlack, false)
	Red     = ColorANSI(ColorSlotRed, false)
	Green   = ColorANSI(ColorSlotGreen, false)
	Yellow  = ColorANSI(ColorSlotYellow, false)
	Blue    = ColorANSI(ColorSlotBlue, false)
	Magenta = ColorANSI(ColorSlotMagenta, false)
	Cyan    = ColorANSI(ColorSlotCyan, false)
	White   = ColorANSI(ColorSlotWhite, false)

	BrightBlack   = ColorANSI(ColorSlotBlack, true)
	BrightRed     = ColorANSI(ColorSlotRed, true)
	BrightGreen   = ColorANSI(ColorSlotGreen, true)
	BrightYellow  = ColorANSI(ColorSlotYellow, true)
	BrightBlue    = ColorANSI(ColorSlotBlue, true)
	BrightMagenta = ColorANSI(ColorSlotMagenta, true)
	BrightCyan    = ColorANSI(ColorSlotCyan, true)
	BrightWhite   = ColorANSI(ColorSlotWhite, true)
)
