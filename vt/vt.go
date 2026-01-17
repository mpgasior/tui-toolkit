package vt

// See: https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
const (
	ESC = "\x1b"
	CSI = ESC + "["
	OSC = ESC + "]"
	ST  = ESC + "\\"
)
