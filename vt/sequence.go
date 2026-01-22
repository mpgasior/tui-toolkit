package vt

import (
	"context"
	"fmt"

	"github.com/nimelo/tui-go/bufiox"
	"github.com/nimelo/tui-go/iox"
)

//go:generate stringer -type=SequenceType
type SequenceType int

const (
	SeqUnknown SequenceType = iota
	SeqUtf8
	SeqControl
	SeqEscape
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

func TODO(ctx context.Context, i *InputBuffer) (next ScanFn, seq Sequence, err error) {
	return nil, seq, fmt.Errorf("not implemented")
}
