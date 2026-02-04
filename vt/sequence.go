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

	"github.com/mpgasior/tui-toolkit/bufiox"
	"github.com/mpgasior/tui-toolkit/iox"
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
	SeqPaste
)

type Sequence struct {
	Data []byte
	Type SequenceType
}

func (s Sequence) Is(types ...SequenceType) bool {
	return slices.Contains(types, s.Type)
}

type InputBuffer struct {
	*bufiox.ContextReader
	BurstTimeout time.Duration
}

func (i *InputBuffer) PeekBurst(ctx context.Context, n int) ([]byte, error) {
	peekCtx, cancel := context.WithTimeout(ctx, i.BurstTimeout)
	defer cancel()

	buf, err := i.PeekContext(peekCtx, n)

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
			_, err = i.DiscardContext(ctx, 1)
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
		if len(buf) == maxCSILen {
			return nil, seq, ErrInvalidSeq
		}
		return nil, seq, err
	}

	if bytes.HasPrefix(buf, []byte(PasteBegin)) {
		return ScanPaste, seq, nil
	}

	size := 2 + idx + 1
	_, _ = i.DiscardContext(ctx, size)
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

	if idx == -1 {
		if len(buf) == maxOSCLen {
			return nil, seq, ErrInvalidSeq
		}
		return nil, seq, err
	}

	size := idx + termLen
	_, _ = i.DiscardContext(ctx, size)
	return nil, Sequence{Data: buf[:size], Type: SeqOSC}, nil
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

	idx := bytes.Index(buf, []byte(ST))
	if idx == -1 {
		if len(buf) == maxDCSLen {
			return nil, seq, ErrInvalidSeq
		}
		return nil, seq, err
	}

	size := idx + 2
	_, _ = i.DiscardContext(ctx, size)
	return nil, Sequence{Data: buf[:size], Type: SeqDCS}, nil
}

func ScanSS3(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	const SS3Len = 3
	buf, err := i.PeekBurst(ctx, SS3Len)
	if len(buf) != SS3Len {
		return nil, seq, err
	}

	if !bytes.HasPrefix(buf, []byte(SS3)) {
		return nil, seq, ErrInvalidSeq
	}

	if !IsCSIFinalByte(buf[2]) {
		return nil, seq, ErrInvalidSeq
	}

	_, _ = i.DiscardContext(ctx, len(buf))
	return nil, Sequence{Data: buf, Type: SeqSS3}, nil
}

func ScanMeta(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	const MetaLen = 2
	buf, err := i.PeekBurst(ctx, MetaLen)
	if len(buf) != MetaLen {
		return nil, seq, err
	}

	r := rune(buf[1])
	if ok := unicode.IsDigit(r) || unicode.IsLetter(r); !ok {
		return nil, seq, ErrInvalidSeq
	}

	if buf[0] != EscByte {
		return nil, seq, ErrInvalidSeq
	}

	_, _ = i.DiscardContext(ctx, len(buf))
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

	_, _ = i.DiscardContext(ctx, size)
	return nil, Sequence{Data: buf[:size], Type: SeqUTF8}, nil
}

func ScanPaste(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	buf, err := i.PeekBurst(ctx, len(PasteBegin))

	if !bytes.HasPrefix(buf, []byte(PasteBegin)) {
		return next, seq, err
	}

	_, _ = i.DiscardContext(ctx, len(buf))

	var content []byte
	for {
		buf, err = i.PeekBurst(ctx, i.Size())
		if len(buf) == 0 {
			return nil, seq, err
		}

		idx := bytes.Index(buf, []byte(PasteEnd))
		if idx == -1 {
			content = append(content, buf...)
			_, _ = i.DiscardContext(ctx, len(buf))
			continue
		}

		content = append(content, buf[:idx]...)
		_, _ = i.DiscardContext(ctx, idx+len([]byte(PasteEnd)))
		return nil, Sequence{Data: content, Type: SeqPaste}, nil
	}
}

func ScanInvalid(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	size := i.Size()
	acc := make([]byte, 0, size)

	for {
		buf, err := i.PeekBurst(ctx, size)
		if len(buf) > 0 {
			_, err = i.DiscardContext(ctx, len(buf))
			acc = append(acc, buf...)
		}

		if err != nil {
			break
		}
	}

	if len(acc) > 0 {
		return nil, Sequence{Data: acc, Type: SeqUnknown}, err
	}

	return nil, seq, err
}
