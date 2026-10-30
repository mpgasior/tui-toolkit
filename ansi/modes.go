package ansi

var (
	BracketedPasteOn  = []byte(CSI + "?2004h")
	BracketedPasteOff = []byte(CSI + "?2004l")
)
