package termx

import (
	"context"
)

type TerminalInput interface {
	MakeRaw() (restore func() error, err error)
	Close() error
	ReadContext(ctx context.Context, p []byte) (n int, err error)
}
