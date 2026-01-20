package iox

import "context"

type ContextReader interface {
	ReadContext(ctx context.Context, p []byte) (n int, err error)
}
