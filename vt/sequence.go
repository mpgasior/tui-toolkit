package vt

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/nimelo/tui-go/bufiox"
	"github.com/nimelo/tui-go/iox"
)

//go:generate stringer -type=SequenceType
type SequenceType int

const (
	SeqUnknown SequenceType = iota
	SeqEscape
	SeqUtf8
	SeqCSI
	SeqOSC
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
	var b byte
	b, err = i.ReadByteContext(ctx)
	if err != nil {
		return nil, seq, err
	}

	if b == 0x1B {
		return nil, Sequence{Data: []byte{b}, Type: SeqEscape}, err
	}

	err = i.UnreadByte()
	return ScanUtf8, seq, err
}

func ScanUtf8(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	var b byte
	b, err = i.ReadByteContext(ctx)
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
		b, err = i.ReadByteContext(ctx)
		if err != nil {
			return nil, seq, err
		}

		data = append(data, b)
	}

	r, rSize := utf8.DecodeRune(data)
	if r == utf8.RuneError && rSize == 1 {
		return nil, Sequence{Data: data, Type: SeqUnknown}, nil
	}

	return nil, Sequence{Data: data, Type: SeqUtf8}, nil
}

func TODO(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	return nil, seq, fmt.Errorf("not implemented")
}
