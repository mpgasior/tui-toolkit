package iox

import "context"

type ContextReader interface {
	Read(ctx context.Context, p []byte) (n int, err error)
}
