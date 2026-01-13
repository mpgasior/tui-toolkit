package utf16x

import (
	"unicode/utf16"
	"unicode/utf8"
)

type RuneScanner struct {
	surrogate rune
}

func (rs *RuneScanner) Write(c uint16) (rune, bool) {
	r := rune(c)

	if rs.surrogate != 0 {
		high := rs.surrogate
		rs.surrogate = 0

		k := utf16.DecodeRune(high, r)
		if k != utf8.RuneError {
			return k, true
		}
	}

	if utf16.IsSurrogate(r) {
		rs.surrogate = r

		return 0, false
	}

	return r, true
}

func (rs *RuneScanner) Reset() {
	rs.surrogate = 0
}
