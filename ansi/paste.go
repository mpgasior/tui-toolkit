package ansi

const (
	BracketedPasteOn  = CSI + "?2004h"
	BracketedPasteOff = CSI + "?2004l"

	BracketedPastePrefix = CSI + "200~"
	BracketedPasteSuffix = CSI + "201~"
)

var (
	BracketedPasteOnBytes  = []byte(BracketedPasteOn)
	BracketedPasteOffBytes = []byte(BracketedPasteOff)

	BracketedPastePrefixBytes = []byte(BracketedPastePrefix)
	BracketedPasteSuffixBytes = []byte(BracketedPasteSuffix)
)
