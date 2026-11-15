package ansi

import (
	"fmt"
	"strconv"
)

const (
	CursorHide = CSI + "?25l"
	CursorShow = CSI + "?25h"

	CursorBlinkOn  = CSI + "?12h"
	CursorBlinkOff = CSI + "?12l"
)

var (
	CursorHideBytes = []byte(CursorHide)
	CursorShowBytes = []byte(CursorShow)

	CursorBlinkOnBytes  = []byte(CursorBlinkOn)
	CursorBlinkOffBytes = []byte(CursorBlinkOff)
)

type CursorStyle uint8

const (
	CursorDefault         CursorStyle = 0
	CursorBlockBlink      CursorStyle = 1
	CursorBlockSteady     CursorStyle = 2
	CursorUnderlineBlink  CursorStyle = 3
	CursorUnderlineSteady CursorStyle = 4
	CursorBarBlink        CursorStyle = 5
	CursorBarSteady       CursorStyle = 6
)

func AppendCursorStyle(dst []byte, s CursorStyle) []byte {
	dst = append(dst, CSI...)
	dst = strconv.AppendInt(dst, int64(s), 10)
	dst = append(dst, ' ')
	dst = append(dst, 'q')
	return dst
}

func AppendCursorColorRGB(dst []byte, hex uint32) []byte {
	dst = append(dst, OSC...)
	dst = append(dst, fmt.Sprintf("12;#%06x", hex&0xFFFFFF)...)
	dst = append(dst, BEL...)
	return dst
}

const CursorResetColorRGB = OSC + "112" + BEL

func AppendResetCursorColorRGB(dst []byte) []byte {
	return append(dst, CursorResetColorRGB...)
}

func AppendCursorMove(dst []byte, x, y int) []byte {
	dst = append(dst, CSI...)
	dst = strconv.AppendInt(dst, int64(y+1), 10)
	dst = append(dst, ';')
	dst = strconv.AppendInt(dst, int64(x+1), 10)
	dst = append(dst, 'H')
	return dst
}

func AppendCursorHide(dst []byte) []byte {
	return append(dst, CursorHide...)
}

func AppendCursorShow(dst []byte) []byte {
	return append(dst, CursorShow...)
}
