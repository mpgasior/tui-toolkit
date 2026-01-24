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
	BurstTimeout time.Duration
}

func (i *InputBuffer) PeekBurst(ctx context.Context, n int) ([]byte, error) {
	peekCtx, cancel := context.WithTimeout(ctx, i.BurstTimeout)
	defer cancel()

	buf, err := i.PeekContext(peekCtx, n)

	if errors.Is(err, context.DeadlineExceeded) && len(buf) > 0 {
		return buf, nil
	}

	return buf, err
}

type ScanFn func(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error)

var ErrInvalidSeq = fmt.Errorf("the sequence is invalid")

const (
	DefaultEscapeTimeout = 20 * time.Millisecond
)

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
			DefaultEscapeTimeout,
		},
	}

	return s
}

func (s *SequenceScanner) SetBurstTimeout(d time.Duration) {
	s.buf.BurstTimeout = d
}

func (s *SequenceScanner) ScanContext(ctx context.Context) bool {
	state := s.initialState

	for state != nil {
		state, s.sequence, s.err = state(ctx, s.buf)

		if errors.Is(s.err, ErrInvalidSeq) {
			state = ScanInvalid
			continue
		}

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
	buf, err := i.PeekContext(ctx, 1)
	if len(buf) != 1 || err != nil {
		return nil, seq, err
	}

	b := buf[0]
	if b == EscByte {
		buf, err = i.PeekBurst(ctx, 2)
		if len(buf) < 2 {
			_, err = i.Discard(1)
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
	const maxCSILen = 64
	buf, err := i.PeekBurst(ctx, maxCSILen)
	if len(buf) == 0 {
		return nil, seq, err
	}

	if !bytes.HasPrefix(buf, []byte(CSI)) {
		return nil, seq, ErrInvalidSeq
	}

	idx := slices.IndexFunc(buf[2:], IsCSIFinalByte)
	if idx == -1 {
		return nil, seq, ErrInvalidSeq
	}

	size := 2 + idx + 1
	if _, err = i.Discard(size); err != nil {
		return nil, seq, err
	}

	return nil, Sequence{Data: buf[:size], Type: SeqCSI}, nil
}

func ScanOSC(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	const maxOSCLen = 8192
	buf, err := i.PeekBurst(ctx, maxOSCLen)
	if len(buf) == 0 {
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

	return nil, seq, ErrInvalidSeq
}

func ScanDCS(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	const maxDCSLen = 8192
	buf, err := i.PeekBurst(ctx, maxDCSLen)
	if len(buf) == 0 {
		return nil, seq, err
	}

	if !bytes.HasPrefix(buf, []byte(DCS)) {
		return nil, seq, ErrInvalidSeq
	}

	if idx := bytes.Index(buf, []byte(ST)); idx != -1 {
		size := idx + 2
		if _, err = i.Discard(size); err != nil {
			return nil, seq, err
		}

		return nil, Sequence{Data: buf[:size], Type: SeqDCS}, nil
	}

	return nil, seq, ErrInvalidSeq
}

func ScanSS3(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	const SS3Len = 3
	buf, err := i.PeekBurst(ctx, SS3Len)
	if len(buf) == 0 {
		return nil, seq, err
	}

	if len(buf) != SS3Len {
		return nil, seq, ErrInvalidSeq
	}

	if !bytes.HasPrefix(buf, []byte(SS3)) {
		return nil, seq, ErrInvalidSeq
	}

	if !IsCSIFinalByte(buf[2]) {
		return nil, seq, ErrInvalidSeq
	}

	if _, err = i.Discard(SS3Len); err != nil {
		return nil, seq, err
	}

	return nil, Sequence{Data: buf, Type: SeqSS3}, nil
}

func ScanMeta(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	const MetaLen = 2
	buf, err := i.PeekBurst(ctx, MetaLen)
	if len(buf) == 0 {
		return nil, seq, err
	}

	if len(buf) != MetaLen {
		return nil, seq, ErrInvalidSeq
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
	buf, err := i.PeekBurst(ctx, utf8.UTFMax)
	if len(buf) == 0 {
		return nil, seq, err
	}

	r, size := utf8.DecodeRune(buf)
	if r == utf8.RuneError && size == 1 {
		return nil, seq, ErrInvalidSeq
	}

	_, err = i.Discard(size)
	return nil, Sequence{Data: buf[:size], Type: SeqUTF8}, nil
}

func ScanInvalid(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	const maxGarbage = 8192

	buf, err := i.PeekBurst(ctx, maxGarbage)

	if len(buf) > 0 {
		_, err = i.Discard(len(buf))
		return nil, Sequence{Data: buf, Type: SeqUnknown}, err
	}

	return nil, seq, err
}
