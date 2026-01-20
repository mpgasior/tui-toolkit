package bufiox

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/nimelo/tui-go/iox"
)

type ContextSplitFunc = func(data []byte, atEOF bool) (advance int, token []byte, err error)

type ErrAmbiguous struct {
	Wait time.Duration
}

func (e *ErrAmbiguous) Error() string {
	return fmt.Sprintf("wait for %v", e.Wait)
}

func IsErrAmbiguous(err error) (*ErrAmbiguous, bool) {
	var errAmbiguous *ErrAmbiguous
	ok := errors.As(err, &errAmbiguous)

	return errAmbiguous, ok
}

type ContextScanner struct {
	reader iox.ContextReader
	split  ContextSplitFunc
	buffer []byte
	token  []byte
	err    error
}

func NewContextScanner(reader iox.ContextReader) *ContextScanner {
	s := &ContextScanner{
		reader: reader,
		split:  bufio.ScanLines,
	}

	return s
}

func (s *ContextScanner) Split(f ContextSplitFunc) {
	s.split = f
}

func (s *ContextScanner) Scan(ctx context.Context) bool {
	for {
		advance, token, err := s.split(s.buffer, false)

		aErr, ok := IsErrAmbiguous(err)
		if ok {
			readCtx, cancel := context.WithTimeout(ctx, aErr.Wait)
			n, rErr := s.readIntoBuffer(readCtx)
			cancel()

			if ctx.Err() != nil {
				s.err = ctx.Err()
				return false
			}

			if errors.Is(rErr, context.DeadlineExceeded) {
				advance, token, err = s.split(s.buffer, true)
			} else if rErr != nil && rErr != io.EOF {
				s.err = rErr
				return false
			} else if n > 0 {
				continue
			}
		}

		if err != nil && aErr != nil {
			s.err = err
			return false
		}

		if advance > 0 {
			s.token = token
			s.buffer = s.buffer[advance:]
			return true
		}

		if _, rErr := s.readIntoBuffer(ctx); rErr != nil {
			if rErr == io.EOF {
				advance, token, err = s.split(s.buffer, true)
				if advance > 0 {
					s.token = token
					s.buffer = s.buffer[advance:]
					return true
				}
			}

			s.err = rErr
			return false
		}
	}
}

func (s *ContextScanner) readIntoBuffer(ctx context.Context) (n int, err error) {
	tmp := make([]byte, 1024)
	n, err = s.reader.ReadContext(ctx, tmp)
	if err == nil {
		s.buffer = append(s.buffer, tmp[:n]...)
	}

	return n, err
}

func (s *ContextScanner) Bytes() []byte {
	return s.token
}

func (s *ContextScanner) Err() error {
	return s.err
}
