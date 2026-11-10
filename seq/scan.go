package seq

import (
	"bytes"
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/ansi"
)

func Scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	if len(data) == 1 && data[0] == ansi.EscByte {
		return 1, data[:1], nil
	}

	if bytes.HasPrefix(data, ansi.BracketedPastePrefixBytes) {
		return ScanBracketedPaste(data, atEOF)
	}

	if bytes.HasPrefix(data, ansi.CSIPrefixBytes) {
		advance, token, err = ScanCSI(data, atEOF)
		if advance > 0 {
			if bytes.HasPrefix(data, ansi.MouseLegacyPrefixBytes) {
				return ScanMouseClickLegacy(data, atEOF)
			}
		}

		return advance, token, err
	}

	if bytes.HasPrefix(data, ansi.SS3PrefixBytes) {
		return ScanSS3(data, atEOF)
	}

	if bytes.HasPrefix(data, ansi.MetaPrefixBytes) {
		return ScanMeta(data, atEOF)
	}

	return ScanUtf8(data, atEOF)
}

func ScanCSI(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !bytes.HasPrefix(data, ansi.CSIPrefixBytes) {
		return errMalformedSeq(data, "CSI")
	}
	for i := len(ansi.CSIPrefixBytes); i < len(data); i += 1 {
		b := data[i]
		if ansi.CSIFinalRange.Contains(b) {
			tokenLen := i + 1
			return tokenLen, data[:tokenLen], nil
		}

		if !ansi.CSIParamRange.Contains(b) && !ansi.CSIIntermediateRange.Contains(b) {
			return errMalformedSeq(data, "CSI")
		}
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func ScanSS3(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !bytes.HasPrefix(data, ansi.SS3PrefixBytes) {
		return errMalformedSeq(data, "SS3")
	}

	finalByteIndex := len(ansi.SS3PrefixBytes)
	if len(data) < finalByteIndex {
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}

	b := data[finalByteIndex]
	tokenLen := finalByteIndex + 1
	if ansi.SS3FinalRange.Contains(b) {
		return tokenLen, data[:tokenLen], nil
	}

	return errMalformedSeq(data, "SS3")
}

func ScanBracketedPaste(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !bytes.HasPrefix(data, ansi.BracketedPastePrefixBytes) {
		return errMalformedSeq(data, "BracketedPaste")
	}

	for i := len(ansi.BracketedPastePrefix); i < len(data); i += 1 {
		tokenLen := i + 1
		if bytes.HasSuffix(data[:tokenLen], ansi.BracketedPasteSuffixBytes) {
			return tokenLen, data[:tokenLen], nil
		}
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func ScanMeta(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !bytes.HasPrefix(data, ansi.MetaPrefixBytes) {
		return errMalformedSeq(data, "Meta")
	}

	finalByteIndex := len(ansi.MetaPrefixBytes)
	if len(data) < finalByteIndex {
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}

	b := data[finalByteIndex]
	tokenLen := finalByteIndex + 1
	if ansi.MetaFinalRange.Contains(b) {
		return tokenLen, data[:tokenLen], nil
	}

	return errMalformedSeq(data, "Meta")
}

func ScanUtf8(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !utf8.FullRune(data) {
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}

	r, size := utf8.DecodeRune(data)
	if r == utf8.RuneError && size == 1 {
		return errMalformedSeq(data, "UTF8")
	}

	return size, data[:size], nil
}
