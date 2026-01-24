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
	reader iox.ContextReader

	buf  []byte
	pos  int
	size int

	err error
}

func NewContextReader(r iox.ContextReader) *ContextReader {
	return &ContextReader{
		reader: r,
		buf:    make([]byte, 4096),
	}
}

func NewContextReaderSize(r iox.ContextReader, size int) *ContextReader {
	return &ContextReader{
		reader: r,
		buf:    make([]byte, size),
	}
}

func (b *ContextReader) Buffered() int {
	return b.size - b.pos
}

func (b *ContextReader) Discard(n int) (discarded int, err error) {
	size := b.Buffered()
	if size <= 0 {
		return 0, nil
	}

	if size >= n {
		discarded = size - n
		b.pos += n

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
		return b.buf[b.pos:b.size], ErrBufferFull
	}

	var err error
	if avail := b.size - b.pos; avail < n {
		n = avail
		err = b.readErr()
		if err != nil {
			err = ErrBufferFull
		}
	}

	return b.buf[b.pos : b.pos+n], nil
}

func (b *ContextReader) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	if b.Buffered() == 0 && len(p) > len(b.buf) {
		return b.reader.ReadContext(ctx, p)
	}

	if b.Buffered() == 0 {
		if b.err != nil {
			return 0, b.readErr()
		}

		n, b.err = b.reader.ReadContext(ctx, b.buf)
		if n == 0 {
			return 0, b.readErr()
		}
		b.size += n
	}

	n = copy(p, b.buf[b.pos:b.size])
	b.pos += n
	return n, nil
}

func (b *ContextReader) fill(ctx context.Context) {
	if b.pos > 0 {
		n := copy(b.buf, b.buf[b.pos:b.size])
		b.size = n
		b.pos = 0
	}

	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		n, err := b.reader.ReadContext(ctx, b.buf[b.size:])
		b.size += n

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
	b.pos, b.size = 0, 0
}

func (b *ContextReader) readErr() error {
	err := b.err
	b.err = nil

	return err
}
