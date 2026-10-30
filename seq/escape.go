package seq

import (
	"context"
	"time"

	"github.com/mpgasior/tui-toolkit/ansi"
)

type ReadContextFunc func(ctx context.Context, p []byte) (n int, err error)

type EscReader struct {
	src        ReadContextFunc
	ctx        context.Context
	escTimeout time.Duration
}

func NewEscReader(ctx context.Context, src ReadContextFunc, escTimeout time.Duration) *EscReader {
	if escTimeout <= 0 {
		escTimeout = 25 * time.Millisecond
	}

	return &EscReader{
		src:        src,
		ctx:        ctx,
		escTimeout: escTimeout,
	}
}

func (r *EscReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	n, err = r.src(r.ctx, p)
	if err != nil {
		return n, err
	}

	if len(p) == n {
		return n, nil
	}

	if n == 1 && p[0] == ansi.ByteEsc {
		ctx, stop := context.WithTimeout(r.ctx, r.escTimeout)
		defer stop()

		n, err = r.src(ctx, p[1:])
		if err != nil {
			if err == ctx.Err() {
				return 1, nil
			}
			return 1, err
		}

		return n + 1, nil
	}

	return n, err
}
