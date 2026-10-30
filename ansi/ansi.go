package ansi

// https://ecma-international.org/wp-content/uploads/ECMA-48_5th_edition_june_1991.pdf

const (
	ESC = "\x1b"
	CSI = ESC + "["
	OSC = ESC + "]"
	DCS = ESC + "P"
	SS3 = ESC + "O"
	ST  = ESC + "\\"
	BEL = "\a"
)

const (
	ByteEsc = 0x1B
)

var (
	BytesCSIPrefix = []byte(CSI)

	CSIFinalRange        = ByteRange{Lo: 0x40, Hi: 0x7E}
	CSIParamRange        = ByteRange{Lo: 0x30, Hi: 0x3F}
	CSIIntermediateRange = ByteRange{Lo: 0x20, Hi: 0x2F}
)
