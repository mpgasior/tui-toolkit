package termx

import (
	"context"
)

type TerminalInput interface {
	MakeRaw() (restore func() error, err error)
	PeekContext(ctx context.Context) (bool, error)
	ReadContext(ctx context.Context, p []byte) (n int, err error)
	Close() error
}
