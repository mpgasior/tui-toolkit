package ansi

import (
	"strconv"
)

type SGRCode = int

const (
	Reset SGRCode = 0

	Bold      SGRCode = 1
	Dim       SGRCode = 2
	Italic    SGRCode = 3
	Underline SGRCode = 4

	FgBlack   SGRCode = 30
	FgRed     SGRCode = 31
	FgGreen   SGRCode = 32
	FgYellow  SGRCode = 33
	FgBlue    SGRCode = 34
	FgMagenta SGRCode = 35
	FgCyan    SGRCode = 36
	FgWhite   SGRCode = 37
	FgDefault SGRCode = 39

	BgBlack   SGRCode = 40
	BgRed     SGRCode = 41
	BgGreen   SGRCode = 42
	BgYellow  SGRCode = 43
	BgBlue    SGRCode = 44
	BgMagenta SGRCode = 45
	BgCyan    SGRCode = 46
	BgWhite   SGRCode = 47
	BgDefault SGRCode = 49

	FgBrightBlack   SGRCode = 90
	FgBrightRed     SGRCode = 91
	FgBrightGreen   SGRCode = 92
	FgBrightYellow  SGRCode = 93
	FgBrightBlue    SGRCode = 94
	FgBrightMagenta SGRCode = 95
	FgBrightCyan    SGRCode = 96
	FgBrightWhite   SGRCode = 97

	BgBrightBlack   SGRCode = 100
	BgBrightRed     SGRCode = 101
	BgBrightGreen   SGRCode = 102
	BgBrightYellow  SGRCode = 103
	BgBrightBlue    SGRCode = 104
	BgBrightMagenta SGRCode = 105
	BgBrightCyan    SGRCode = 106
	BgBrightWhite   SGRCode = 107
)

const (
	SGRStart    = CSI
	SGREnd      = 'm'
	SGRParamSep = ';'
)

func AppendSGRStart(dst []byte) []byte {
	return append(dst, SGRStart...)
}

func AppendSGREnd(dst []byte) []byte {
	return append(dst, SGREnd)
}

func AppendSGRReset(dst []byte) []byte {
	dst = append(dst, SGRStart...)
	dst = append(dst, byte(Reset))
	return append(dst, SGREnd)
}

func AppendSGRParam(dst []byte, code int) []byte {
	if len(dst) > 0 {
		last := dst[len(dst)-1]
		if last != ';' && last != '[' {
			dst = append(dst, ';')
		}
	}

	return strconv.AppendInt(dst, int64(code), 10)
}

func AppendSGRRGBFg(dst []byte, hex uint32) []byte {
	if len(dst) > 0 {
		last := dst[len(dst)-1]
		if last != ';' && last != '[' {
			dst = append(dst, ';')
		}
	}

	mask := uint32(0xFF)
	r := int64((hex >> 16) & mask)
	g := int64((hex >> 8) & mask)
	b := int64(hex & mask)

	dst = append(dst, "38;2;"...)
	dst = strconv.AppendInt(dst, r, 10)
	dst = append(dst, ';')
	dst = strconv.AppendInt(dst, g, 10)
	dst = append(dst, ';')
	dst = strconv.AppendInt(dst, b, 10)
	return dst
}

func AppendSGRRGBBg(dst []byte, hex uint32) []byte {
	if len(dst) > 0 {
		last := dst[len(dst)-1]
		if last != ';' && last != '[' {
			dst = append(dst, ';')
		}
	}

	mask := uint32(0xFF)
	r := int64((hex >> 16) & mask)
	g := int64((hex >> 8) & mask)
	b := int64(hex & mask)

	dst = append(dst, "48;2;"...)
	dst = strconv.AppendInt(dst, r, 10)
	dst = append(dst, ';')
	dst = strconv.AppendInt(dst, g, 10)
	dst = append(dst, ';')
	dst = strconv.AppendInt(dst, b, 10)
	return dst
}
