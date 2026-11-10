package ansi

const (
	AltScreenOn  = CSI + "?1049h"
	AltScreenOff = CSI + "?1049l"
)

var (
	AltScreenOnBytes  = []byte(AltScreenOn)
	AltScreenOffBytes = []byte(AltScreenOff)
)
