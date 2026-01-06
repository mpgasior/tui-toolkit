package vt

// Simple Cursor Positioning
const (
	ReverseIndex  = ESC + "M"
	SaveCursor    = ESC + "7"
	RestoreCursor = ESC + "8"
)

// Cursor Positioning
const (
	CursorUpFmt                        = CSI + "%dA"
	CursorDownFmt                      = CSI + "%dB"
	CursorRightFmt                     = CSI + "%dC"
	CursorLeftFmt                      = CSI + "%dD"
	CursorNextLineFmt                  = CSI + "%dE"
	CursorPrevLineFmt                  = CSI + "%dF"
	CursorHorizontalLineFmt            = CSI + "%dG"
	CursorVerticalLineFmt              = CSI + "%dd"
	CursorPositionFmt                  = CSI + "%d;%dH"
	CursorHorizontalVerticalPostionFmt = CSI + "%d;%df"
	SaveCursorAnsi                     = CSI + "s"
	RestoreCursorAnsi                  = CSI + "u"
)

// Cursor Visibility
const (
	CursorBlink   = CSI + "?12h"
	CursorNoBlink = CSI + "?12l"
	CursorShow    = CSI + "?25h"
	CursorHide    = CSI + "?25l"
)

// Cursor Shape
const (
	CursorShapeDefault        = CSI + "0 q"
	CursorShapeBlockBlink     = CSI + "1 q"
	CursorShapeBlock          = CSI + "2 q"
	CursorShapeUnderlineBlink = CSI + "3 q"
	CursorShapeUnderline      = CSI + "4 q"
	CursorShapeBarBlink       = CSI + "5 q"
	CursorShapeBar            = CSI + "6 q"
)
