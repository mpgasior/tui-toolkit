package hook

import "github.com/mpgasior/tui-toolkit/ansi"

type KeyEvent struct {
	R       rune
	Code    ansi.Key
	stopped bool
}

func (ke *KeyEvent) StopPropagation() {
	ke.stopped = true
}
