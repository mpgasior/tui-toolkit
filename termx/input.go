package termx

import (
	"context"
)

type TerminalInput interface {
	MakeRaw() (restore func() error, err error)
	Ready(ctx context.Context) error
	Read(ctx context.Context, p []byte) (n int, err error)
	Close() error
}
