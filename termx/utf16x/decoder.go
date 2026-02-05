package utf16x

import (
	"unicode/utf16"
	"unicode/utf8"
)

type Decoder struct {
	lastSurrogate rune
}

func (d *Decoder) Decode(u uint16) (rune, bool) {
	r := rune(u)

	if d.lastSurrogate != 0 {
		high := d.lastSurrogate
		d.lastSurrogate = 0

		combined := utf16.DecodeRune(high, r)
		if combined != utf8.RuneError {
			return combined, true
		}
	}

	if utf16.IsSurrogate(r) {
		d.lastSurrogate = r

		return 0, false
	}

	return r, true
}

func (d *Decoder) Reset() {
	d.lastSurrogate = 0
}
