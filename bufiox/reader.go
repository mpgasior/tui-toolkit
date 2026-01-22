package bufiox

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/nimelo/tui-go/iox"
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

func (b *ContextReader) Buffered() int {
	return b.size - b.pos
}

func (b *ContextReader) Discard(n int) (discarded int, err error) {
	size := b.Buffered()
	if size <= 0 {
		return 0, fmt.Errorf("empty buffer")
	}

	if size >= n {
		discarded = size - n
		b.pos += n

		return discarded, nil
	}

	discarded = n - b.Buffered()
	return discarded, fmt.Errorf("not enough in the buffer")
}

func (b *ContextReader) PeekContext(ctx context.Context, n int) ([]byte, error) {
	for b.Buffered() < n && b.err == nil {
		b.fill(ctx)
	}

	if b.Buffered() < n {
		return b.buf[b.pos:b.size], b.readErr()
	}

	return b.buf[b.pos : b.pos+n], nil
}

func (b *ContextReader) ReadContext(ctx context.Context, p []byte) (int, error) {
	if b.pos >= b.size && b.err == nil {
		b.fill(ctx)
	}

	if b.pos >= b.size {
		return 0, b.readErr()
	}

	n := copy(p, b.buf[b.pos:b.size])
	b.pos += n
	return n, nil
}

func (b *ContextReader) fill(ctx context.Context) {
	if b.pos > 0 {
		n := copy(b.buf, b.buf[b.pos:b.size])
		b.size = n
		b.pos = 0
	}

	if b.size == len(b.buf) {
		tmp := make([]byte, b.size*2)
		copy(tmp, b.buf[b.pos:b.size])
		b.buf = tmp
	}

	n, err := b.reader.ReadContext(ctx, b.buf[b.size:])
	b.size += n

	if err != nil {
		b.err = err
	}
}

func (b *ContextReader) ReadByteContext(ctx context.Context) (byte, error) {
	if b.Buffered() < 1 {
		b.fill(ctx)
	}

	if b.Buffered() >= 1 {
		byte := b.buf[b.pos]
		b.pos += 1

		return byte, nil
	}

	return 0, b.readErr()
}

func (b *ContextReader) ReadRuneContext(ctx context.Context) (r rune, size int, err error) {
	for {
		if b.Buffered() > 0 {
			buf := b.buf[b.pos:b.size]
			r, size = utf8.DecodeRune(buf)
			if utf8.FullRune(buf) {
				b.pos += size
				return r, size, nil
			}
		}

		if b.err != nil {
			if b.Buffered() <= 0 {
				return 0, 0, b.readErr()
			}

			r, size = utf8.DecodeRune(b.buf[b.pos:b.size])
			b.pos += size
			return r, size, nil
		}

		b.fill(ctx)
	}
}

func (b *ContextReader) Reset() {
	b.pos, b.size = 0, 0
}

func (b *ContextReader) UnreadByte() error {
	if b.pos > 0 {
		b.pos -= 1
		return nil
	}

	return fmt.Errorf("nothing to unread")
}

func (b *ContextReader) readErr() error {
	err := b.err
	b.err = nil

	return err
}
