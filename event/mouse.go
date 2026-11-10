package event

import (
	"bytes"

	"github.com/mpgasior/tui-toolkit/ansi"
)

type Mouse struct {
	X, Y byte
}

type MouseDecoder struct {
}

func NewMouseDecoder() *MouseDecoder {
	return &MouseDecoder{}
}

func (d *MouseDecoder) Decode(buf []byte) (mouse Mouse, ok bool) {
	if legacy, ok := d.decodeMouseLegacy(buf); ok {
		return legacy, ok
	}

	return mouse, true
}

func (d *MouseDecoder) decodeMouseLegacy(buf []byte) (mouse Mouse, ok bool) {
	if !bytes.HasPrefix(buf, ansi.MouseLegacyPrefixBytes) {
		return mouse, false
	}

	if len(buf) != ansi.MouseLegacyLength {
		return mouse, false
	}

	x := buf[ansi.MouseLegacyColIndex] - ansi.MouseLegacyColRowOffset
	y := buf[ansi.MouseLegacyRowIndex] - ansi.MouseLegacyColRowOffset

	mouse = Mouse{
		X: x, Y: y,
	}

	return mouse, true
}
