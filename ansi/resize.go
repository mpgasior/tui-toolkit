package ansi

const (
	InBandResizeModeOn  = CSI + "?2048h"
	InBandResizeModeOff = CSI + "?2048l"
)

var (
	InBandResizeModeOnBytes  = []byte(InBandResizeModeOn)
	InBandResizeModeOffBytes = []byte(InBandResizeModeOff)
)

const InBandResizePrefix = CSI + "48;"

var (
	InBandResizePrefixBytes = []byte(InBandResizePrefix)
	InBandResizeFinalRange  = ByteRange{Lo: 't', Hi: 't'}
)
