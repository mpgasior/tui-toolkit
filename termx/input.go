package termx

import "github.com/nimelo/tui-go/iox"

type TerminalInput interface {
	iox.ContextReader
	MakeRaw() (restore func() error, err error)
	Close() error
}
