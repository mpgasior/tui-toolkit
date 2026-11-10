package seq

import (
	"bytes"

	"github.com/mpgasior/tui-toolkit/ansi"
)

func ScanMouseClickLegacy(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !bytes.HasPrefix(data, ansi.MouseLegacyPrefixBytes) {
		return errMalformedSeq(data, "CSI Mouse Legacy")
	}

	if len(data) < ansi.MouseLegacyLength {
		if atEOF {
			return errMalformedSeq(data, "CSI Mouse Legacy")
		}
		return 0, nil, nil
	}

	tokenLen := ansi.MouseLegacyLength
	return tokenLen, data[:tokenLen], nil
}
