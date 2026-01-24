package bufiox

import (
	"context"
	"errors"
	"io"

	"github.com/nimelo/tui-go/iox"
)

var (
	ErrBufferFull = errors.New("bufiox: buffer full")
)

const (
	maxConsecutiveEmptyReads = 100
)

type ContextReader struct {
	rd  iox.ContextReader
	buf []byte
	r   int
	w   int
	err error
}

func NewContextReader(r iox.ContextReader) *ContextReader {
	return &ContextReader{
		rd:  r,
		buf: make([]byte, 4096),
	}
}

func NewContextReaderSize(r iox.ContextReader, size int) *ContextReader {
	return &ContextReader{
		rd:  r,
		buf: make([]byte, size),
	}
}

func (b *ContextReader) Buffered() int {
	return b.w - b.r
}

func (b *ContextReader) Discard(n int) (discarded int, err error) {
	size := b.Buffered()
	if size <= 0 {
		return 0, nil
	}

	if size >= n {
		discarded = size - n
		b.r += n

		return discarded, nil
	}

	discarded = n - b.Buffered()
	return discarded, b.readErr()
}

func (b *ContextReader) PeekContext(ctx context.Context, n int) ([]byte, error) {
	if n > len(b.buf) {
		return nil, ErrBufferFull
	}

	for b.Buffered() < n && b.Buffered() < len(b.buf) && b.err == nil {
		b.fill(ctx)
	}

	if n > len(b.buf) {
		return b.buf[b.r:b.w], ErrBufferFull
	}

	if b.Buffered() < n {
		return b.buf[b.r:b.w], b.readErr()
	}

	return b.buf[b.r : b.r+n], nil
}

func (b *ContextReader) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	if b.Buffered() == 0 && len(p) > len(b.buf) {
		return b.rd.ReadContext(ctx, p)
	}

	if b.Buffered() == 0 {
		if b.err != nil {
			return 0, b.readErr()
		}

		n, b.err = b.rd.ReadContext(ctx, b.buf)
		if n == 0 {
			return 0, b.readErr()
		}
		b.w += n
	}

	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	return n, nil
}

func (b *ContextReader) fill(ctx context.Context) {
	if b.r > 0 {
		n := copy(b.buf, b.buf[b.r:b.w])
		b.w = n
		b.r = 0
	}

	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		n, err := b.rd.ReadContext(ctx, b.buf[b.w:])
		b.w += n

		if err != nil {
			b.err = err
			return
		}

		if n > 0 {
			return
		}
	}

	b.err = io.ErrNoProgress
}

func (b *ContextReader) Reset() {
	b.r, b.w = 0, 0
}

func (b *ContextReader) readErr() error {
	err := b.err
	b.err = nil

	return err
}
