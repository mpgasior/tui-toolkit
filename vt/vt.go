package vt

// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
const (
	ESC = "\x1b"
	CSI = ESC + "["
	OSC = ESC + "]"
	ST  = ESC + "\\"
	SS3 = ESC + "O"
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
