package termio

import "context"

type Reader interface {
	ReadContext(ctx context.Context, p []byte) (n int, err error)
	Close() error
}
