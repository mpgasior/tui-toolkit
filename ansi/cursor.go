package ansi

import "strconv"

const (
	CursorHide = CSI + "?25l"
	CursorShow = CSI + "?25h"

	CursorBlinkOn  = CSI + "?12h"
	CursorBlinkOff = CSI + "?12l"
)

func AppendCursorMove(dst []byte, x, y int) []byte {
	dst = append(dst, CSI...)
	dst = strconv.AppendInt(dst, int64(y), 10)
	dst = append(dst, ';')
	dst = strconv.AppendInt(dst, int64(x), 10)
	dst = append(dst, 'f')
	return dst
}

func AppendCursorHide(dst []byte) []byte {
	return append(dst, CursorHide...)
}

func AppendCursorShow(dst []byte) []byte {
	return append(dst, CursorShow...)
}
