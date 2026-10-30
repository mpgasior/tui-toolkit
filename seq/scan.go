package seq

import (
	"bytes"
	"errors"
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/ansi"
)

var ErrMalformedSequence = errors.New("seq: malformed ansi escape sequence")

func Scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	if bytes.HasPrefix(data, ansi.BytesBracketedPastePrefix) {
		return ScanBracketedPaste(data, atEOF)
	}

	if bytes.HasPrefix(data, ansi.BytesCSIPrefix) {
		return ScanCSI(data, atEOF)
	}

	if bytes.HasPrefix(data, ansi.BytesSS3Prefix) {
		return ScanSS3(data, atEOF)
	}

	return ScanUtf8(data, atEOF)
}

func ScanCSI(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !bytes.HasPrefix(data, ansi.BytesCSIPrefix) {
		return len(data), data, ErrMalformedSequence
	}
	for i := len(ansi.BytesCSIPrefix); i < len(data); i += 1 {
		b := data[i]
		if ansi.CSIFinalRange.Contains(b) {
			tokenLen := i + 1
			return tokenLen, data[:tokenLen], nil
		}

		if !ansi.CSIParamRange.Contains(b) && !ansi.CSIIntermediateRange.Contains(b) {
			return i, data[:i], ErrMalformedSequence
		}
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func ScanSS3(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !bytes.HasPrefix(data, ansi.BytesSS3Prefix) {
		return len(data), data, ErrMalformedSequence
	}

	finalByteIndex := len(ansi.BytesSS3Prefix)
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

	return tokenLen, data[:tokenLen], ErrMalformedSequence
}

func ScanBracketedPaste(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !bytes.HasPrefix(data, ansi.BytesBracketedPastePrefix) {
		return len(data), data, ErrMalformedSequence
	}

	for i := len(ansi.BytesBracketedPastePrefix); i < len(data); i += 1 {
		tokenLen := i + 1
		if bytes.HasSuffix(data[:tokenLen], ansi.BytesBracketedPasteSuffix) {
			return tokenLen, data[:tokenLen], nil
		}
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func ScanUtf8(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	if !utf8.FullRune(data) {
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}

	r, size := utf8.DecodeRune(data)
	if r == utf8.RuneError && size == 1 {
		return 0, nil, ErrMalformedSequence
	}

	return size, data[:size], nil
}
