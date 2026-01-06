package vt

//go:generate go run golang.org/x/tools/cmd/stringer -type=Key,SequenceType,Attr,FgColor,FgBrightColor,BgColor,BgBrightColor -output=vt_string.go

// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
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
	EscByte byte = 0x1B
)

func IsESCFinalByte(b byte) bool {
	return b >= 0x30 && b <= 0x7E
}

func IsCSIFinalByte(b byte) bool {
	return b >= 0x40 && b <= 0x7E
}
