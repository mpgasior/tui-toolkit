package termx

import "github.com/mpgasior/tui-go/iox"

type TerminalInput interface {
	iox.ContextReader
	MakeRaw() (restore func() error, err error)
	Close() error
}
