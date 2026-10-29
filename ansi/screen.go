package ansi

var (
	AltScreenOn  = []byte("\x1b[?1049h")
	AltScreenOff = []byte("\x1b[?1049l")
)
