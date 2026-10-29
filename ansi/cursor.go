package ansi

var (
	CursorHide = []byte("\x1b[?25l")
	CursorShow = []byte("\x1b[?25h")

	CursorBlinkOn  = []byte("\x1b[?12h")
	CursorBlinkOff = []byte("\x1b[?12l")
)
