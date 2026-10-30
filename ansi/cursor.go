package ansi

var (
	CursorHide = []byte(CSI + "?25l")
	CursorShow = []byte(CSI + "?25h")

	CursorBlinkOn  = []byte(CSI + "?12h")
	CursorBlinkOff = []byte(CSI + "?12l")
)
