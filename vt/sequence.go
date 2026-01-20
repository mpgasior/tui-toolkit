package vt

import (
	"bytes"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/nimelo/tui-go/bufiox"
)

var ErrInvalidSequence = fmt.Errorf("invalid sequence")

func ScanMulti(funcs ...bufiox.ContextSplitFunc) bufiox.ContextSplitFunc {
	return nil
}

func ScanInput(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = ScanESC(data, atEOF)

	if token != nil {
		return advance, token, err
	}

	if err != nil {
		if _, ok := bufiox.IsErrAmbiguous(err); ok {
			return advance, token, err
		}

		advance, token, err = ScanSS3(data, atEOF)
	}

	return advance, token, err
}

func ScanESC(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		if atEOF {
			return 0, nil, ErrInvalidSequence
		}

		return 0, nil, nil
	}

	if data[0] == 0x1b {
		if atEOF {
			return 1, data[0:0], nil
		}

		if len(data) == 1 {
			err := &bufiox.ErrAmbiguous{
				Wait: 20 * time.Millisecond,
			}
			return 0, nil, err
		}
	}

	return 0, nil, ErrInvalidSequence
}

func ScanCSI(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < len(CSI) {
		if atEOF && len(data) > 0 {
			return 0, nil, ErrInvalidSequence
		}

		return 0, nil, nil
	}

	if bytes.HasPrefix(data, []byte(CSI)) {
		return 0, nil, ErrInvalidSequence
	}

	for idx, b := range data[len(CSI):] {
		if IsFinalByte(b) {
			advance := idx + 1 + len(CSI)
			return advance, data[:advance], nil
		}
	}

	if atEOF {
		return 0, nil, ErrInvalidSequence
	}

	return 0, nil, nil
}

func ScanSS3(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < 3 {
		if atEOF && len(data) > 0 {
			return 0, nil, ErrInvalidSequence
		}

		return 0, nil, nil
	}

	if !bytes.HasPrefix(data, []byte(SS3)) {
		return 0, nil, ErrInvalidSequence
	}

	finalByte := data[2]

	if IsFinalByte(finalByte) {
		return 3, data[:3], nil
	}

	return 0, nil, ErrInvalidSequence
}

func ScanUtf8(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	r, size := utf8.DecodeRune(data)

	if r == utf8.RuneError && size == 1 {
		if len(data) < utf8.UTFMax && !atEOF {
			return 0, nil, nil
		}

		return 0, nil, ErrInvalidSequence
	}

	return size, data[:size], nil
}

func IsFinalByte(b byte) bool {
	if b >= 0x40 && b <= 0x7E {
		return true
	}

	return false
}
