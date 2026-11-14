package event

import (
	"bytes"

	"github.com/mpgasior/tui-toolkit/ansi"
)

type Paste []byte

type PasteDecoder struct{}

func NewPasteDecoder() *PasteDecoder {
	return &PasteDecoder{}
}

func (d *PasteDecoder) Decode(buf []byte) (paste Paste, ok bool) {
	if !bytes.HasPrefix(buf, ansi.BracketedPastePrefixBytes) {
		return paste, ok
	}

	if !bytes.HasSuffix(buf, ansi.BracketedPasteSuffixBytes) {
		return paste, ok
	}

	start, end := len(ansi.BracketedPastePrefixBytes), len(buf)-len(ansi.BracketedPasteSuffixBytes)
	return buf[start:end], true
}
