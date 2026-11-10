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
	EscByte = 0x1B
)

const CSIPrefix = CSI

var (
	CSIPrefixBytes = []byte(CSI)

	CSIFinalRange        = ByteRange{Lo: 0x40, Hi: 0x7E}
	CSIParamRange        = ByteRange{Lo: 0x30, Hi: 0x3F}
	CSIIntermediateRange = ByteRange{Lo: 0x20, Hi: 0x2F}
)

const (
	OSCPrefix  = OSC
	OSCSuffix  = ST
	OSCSuffix2 = ESC + BEL
)

var (
	OSCPrefixBytes  = []byte(OSCPrefix)
	OSCSuffixBytes  = []byte(OSCSuffix)
	OSCSuffix2Bytes = []byte(OSCSuffix2)

	OSCIntermediateRange = RangeTable{
		{Lo: 0x88, Hi: 0x8D},
		{Lo: 0xA0, Hi: 0xFE},
	}
)

const SS3Prefix = SS3

var (
	SS3PrefixBytes = []byte(SS3Prefix)

	SS3FinalRange = ByteRange{Lo: 0x40, Hi: 0x7E}
)

const MetaPrefix = ESC

var (
	MetaPrefixBytes = []byte(MetaPrefix)

	MetaFinalRange = RangeTable{
		{Lo: 0x00, Hi: 0x1F},
		{Lo: 0x7E, Hi: 0x7E},
	}
)
