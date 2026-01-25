package bufiox

import (
	"context"
	"errors"
	"io"

	"github.com/mpgasior/tui-go/iox"
)

const (
	defaultBufSize = 4096
)

var (
	ErrBufferFull    = errors.New("bufiox: buffer full")
	ErrNegativeCount = errors.New("bufiox: negative count")
)

type ContextReader struct {
	rd  iox.ContextReader
	buf []byte
	r   int
	w   int
	err error
}

const minReadBufferSize = 16
const maxConsecutiveEmptyReads = 100

func NewContextReaderSize(r iox.ContextReader, size int) *ContextReader {
	return &ContextReader{
		rd:  r,
		buf: make([]byte, max(minReadBufferSize, size)),
	}
}

func NewContextReader(r iox.ContextReader) *ContextReader {
	return &ContextReader{
		rd:  r,
		buf: make([]byte, defaultBufSize),
	}
}

func (b *ContextReader) Size() int {
	return len(b.buf)
}

func (b *ContextReader) Reset(r iox.ContextReader) {
	if r == b {
		return
	}

	if b.buf == nil {
		b.buf = make([]byte, defaultBufSize)
	}

	b.reset(b.buf, r)
}

func (b *ContextReader) reset(buf []byte, r iox.ContextReader) {
	*b = ContextReader{
		buf: buf,
		rd:  r,
	}
}

var errNegativeRead = errors.New("bufiox: reader returned negative count from Read")

func (b *ContextReader) fill(ctx context.Context) {
	if b.r > 0 {
		n := copy(b.buf, b.buf[b.r:b.w])
		b.r, b.w = 0, n
	}

	for i := maxConsecutiveEmptyReads; i > 0; i-- {
		n, err := b.rd.ReadContext(ctx, b.buf[b.w:])
		if n < 0 {
			panic(errNegativeRead)
		}

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

func (b *ContextReader) readErr() error {
	err := b.err
	b.err = nil

	return err
}

func (b *ContextReader) PeekContext(ctx context.Context, n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrNegativeCount
	}

	for b.w-b.r < n && b.w-b.r < len(b.buf) && b.err == nil {
		b.fill(ctx)
	}

	err := b.readErr()
	if n > len(b.buf) {
		n = len(b.buf)
		if err == nil {
			err = ErrBufferFull
		}
	}

	if avail := b.w - b.r; avail < n {
		n = avail
		if err == nil {
			err = ErrBufferFull
		}
	}

	return b.buf[b.r : b.r+n], err
}

func (b *ContextReader) DiscardContext(ctx context.Context, n int) (discarded int, err error) {
	if n < 0 {
		return 0, ErrNegativeCount
	}

	if n == 0 {
		return
	}

	remain := n
	for {
		skip := b.Buffered()
		if skip == 0 {
			b.fill(ctx)
			skip = b.Buffered()
		}

		if skip > remain {
			skip = remain
		}

		b.r += skip
		remain -= skip
		if remain == 0 {
			return n, b.readErr()
		}

		if b.err != nil {
			return n - remain, b.readErr()
		}
	}
}

func (b *ContextReader) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	n = len(p)

	if n == 0 {
		if b.Buffered() > 0 {
			return 0, nil
		}
		return 0, b.readErr()
	}

	if b.r == b.w {
		if b.err != nil {
			return 0, b.readErr()
		}
		if len(p) >= len(b.buf) {
			n, b.err = b.rd.ReadContext(ctx, p)
			if n < 0 {
				panic(errNegativeRead)
			}

			return n, b.readErr()
		}

		b.r = 0
		b.w = 0
		n, b.err = b.rd.ReadContext(ctx, b.buf)
		if n < 0 {
			panic(errNegativeRead)
		}

		b.w += n
	}

	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	return n, b.readErr()
}

func (b *ContextReader) Buffered() int {
	return b.w - b.r
}
