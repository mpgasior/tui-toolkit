package vt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/nimelo/tui-go/bufiox"
	"github.com/nimelo/tui-go/iox"
)

type SequenceType int

const (
	SeqUnknown SequenceType = iota
	SeqEscape
	SeqUTF8
	SeqCSI
	SeqOSC
	SeqDCS
	SeqSS3
	SeqMeta
)

type Sequence struct {
	Data []byte
	Type SequenceType
}

func (s Sequence) Is(t SequenceType) bool {
	return s.Type == t
}

type InputBuffer struct {
	*bufiox.ContextReader
}

type ScanFn func(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error)

var ErrInvalidSeq = fmt.Errorf("the sequence is invalid")

type SequenceScanner struct {
	buf          *InputBuffer
	initialState ScanFn
	sequence     Sequence
	err          error
}

func NewSequenceScanner(r iox.ContextReader, f ScanFn) *SequenceScanner {
	s := &SequenceScanner{
		initialState: f,
		buf: &InputBuffer{
			bufiox.NewContextReader(r),
		},
	}

	return s
}

func (s *SequenceScanner) ScanContext(ctx context.Context) bool {
	state := s.initialState

	for state != nil {
		state, s.sequence, s.err = state(ctx, s.buf)

		if s.err != nil {
			return false
		}

		if len(s.sequence.Data) > 0 {
			return true
		}
	}

	return false
}

func (s *SequenceScanner) Sequence() Sequence {
	return s.sequence
}

func (s *SequenceScanner) Err() error {
	return s.err
}

func ScanInitial(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	var buf []byte
	buf, err = i.PeekContext(ctx, 1)
	if len(buf) != 1 || err != nil {
		return nil, seq, err
	}

	b := buf[0]
	if b == EscByte {
		peekCtx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
		defer cancel()

		buf, err = i.PeekContext(peekCtx, 2)
		if len(buf) < 2 {
			_, err = i.ReadByteContext(ctx)
			return nil, Sequence{Data: []byte{b}, Type: SeqEscape}, err
		}

		prefix := buf[:2]
		if bytes.HasPrefix(prefix, []byte(CSI)) {
			return ScanCSI, seq, err
		}

		if bytes.HasPrefix(prefix, []byte(OSC)) {
			return ScanOSC, seq, err
		}

		if bytes.HasPrefix(prefix, []byte(DCS)) {
			return ScanDCS, seq, err
		}

		if bytes.HasPrefix(prefix, []byte(SS3)) {
			return ScanSS3, seq, err
		}

		if next, seq, err = ScanMeta(ctx, i); err != nil {
			return next, seq, err
		}
	}

	return ScanUtf8, seq, err
}

func ScanCSI(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	var buf []byte

	peekCtx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()

	const maxCSILen = 64
	buf, err = i.PeekContext(peekCtx, maxCSILen)
	if len(buf) == 0 {
		return nil, seq, err
	}

	if !bytes.HasPrefix(buf, []byte(CSI)) {
		return nil, seq, ErrInvalidSeq
	}

	idx := slices.IndexFunc(buf[2:], IsCSIFinalByte)
	if idx == -1 {
		if errors.Is(err, context.DeadlineExceeded) && len(buf) <= maxCSILen {
			return nil, seq, err
		}

		return nil, seq, ErrInvalidSeq
	}

	size := 2 + idx + 1
	if _, err = i.Discard(size); err != nil {
		return nil, seq, err
	}

	return nil, Sequence{Data: buf[:size], Type: SeqCSI}, nil
}

func ScanOSC(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	peekCtx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()

	const maxOSCLen = 8192
	buf, err := i.PeekContext(peekCtx, maxOSCLen)
	if len(buf) < 2 {
		return nil, seq, err
	}

	if !bytes.HasPrefix(buf, []byte(OSC)) {
		return nil, seq, ErrInvalidSeq
	}

	idx := bytes.Index(buf, []byte(ST))
	termLen := 2

	if idx == -1 {
		idx = bytes.Index(buf, []byte(ESC+BEL))
		termLen = 1
	}

	if idx != -1 {
		size := idx + termLen
		if _, err = i.Discard(size); err != nil {
			return nil, seq, err
		}

		return nil, Sequence{Data: buf[:size], Type: SeqOSC}, nil
	}

	if len(buf) >= maxOSCLen {
		return nil, seq, ErrInvalidSeq
	}

	return nil, seq, err
}

func ScanDCS(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	peekCtx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()

	const maxDCSLen = 8192
	buf, err := i.PeekContext(peekCtx, maxDCSLen)
	if len(buf) < 2 {
		return nil, seq, err
	}

	if !bytes.HasPrefix(buf, []byte(DCS)) {
		return nil, seq, ErrInvalidSeq
	}

	idx := bytes.Index(buf, []byte(ST))

	if idx != -1 {
		size := idx + 2
		if _, err = i.Discard(size); err != nil {
			return nil, seq, err
		}

		return nil, Sequence{Data: buf[:size], Type: SeqDCS}, nil
	}

	if len(buf) >= maxDCSLen {
		return nil, seq, ErrInvalidSeq
	}

	return nil, seq, err
}

func ScanSS3(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	peekCtx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()

	const SS3Len = 3
	buf, err := i.PeekContext(peekCtx, SS3Len)
	if err != nil || len(buf) != SS3Len {
		return nil, seq, err
	}

	if bytes.HasPrefix(buf, []byte(SS3)) && IsCSIFinalByte(buf[2]) {
		_, err = i.ReadContext(ctx, buf)
		return nil, Sequence{Data: buf, Type: SeqSS3}, nil
	}

	return nil, seq, ErrInvalidSeq
}

func ScanMeta(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	peekCtx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()

	const MetaLen = 2
	buf, err := i.PeekContext(peekCtx, MetaLen)
	if err != nil || len(buf) != MetaLen {
		return nil, seq, err
	}

	r := rune(buf[1])
	if ok := unicode.IsDigit(r) || unicode.IsLetter(r); !ok {
		return nil, seq, ErrInvalidSeq
	}

	if buf[0] != EscByte {
		return nil, seq, ErrInvalidSeq
	}

	if _, err := i.Discard(MetaLen); err != nil {
		return nil, seq, err
	}

	return nil, Sequence{Data: buf, Type: SeqMeta}, nil
}

func ScanUtf8(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	shortCtx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()

	var b byte
	b, err = i.ReadByteContext(shortCtx)
	if err != nil {
		return ScanInitial, seq, err
	}

	if !utf8.RuneStart(b) {
		return nil, Sequence{Data: []byte{b}, Type: SeqUnknown}, nil
	}

	size := utf8.RuneLen(rune(b))
	if size == -1 {
		return nil, Sequence{Data: []byte{b}, Type: SeqUnknown}, nil
	}

	data := make([]byte, 1, utf8.UTFMax)
	data[0] = b

	for !utf8.FullRune(data) {
		b, err = i.ReadByteContext(shortCtx)
		if err != nil {
			return nil, seq, err
		}

		data = append(data, b)
	}

	r, rSize := utf8.DecodeRune(data)
	if r == utf8.RuneError && rSize == 1 {
		return nil, Sequence{Data: data, Type: SeqUnknown}, nil
	}

	return nil, Sequence{Data: data, Type: SeqUTF8}, nil
}

func TODO(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	return nil, seq, fmt.Errorf("not implemented")
}
