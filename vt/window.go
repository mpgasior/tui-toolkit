package vt

const (
	SetWindowTitleFmt = OSC + "2;%s" + ST

	UseAltScreenBuffer  = CSI + "?1049h"
	UseMainScreenBuffer = CSI + "?1049l"
)
