package utf16x

import (
	"unicode/utf16"
	"unicode/utf8"
)

type RuneScanner struct {
	lastSurrogate rune
}

func (rs *RuneScanner) Scan(unit uint16) (rune, bool) {
	r := rune(unit)

	if rs.lastSurrogate != 0 {
		high := rs.lastSurrogate
		rs.lastSurrogate = 0

		k := utf16.DecodeRune(high, r)
		if k != utf8.RuneError {
			return k, true
		}
	}

	if utf16.IsSurrogate(r) {
		rs.lastSurrogate = r

		return 0, false
	}

	return r, true
}

func (rs *RuneScanner) Reset() {
	rs.lastSurrogate = 0
}
