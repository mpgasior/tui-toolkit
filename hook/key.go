package hook

import "github.com/mpgasior/tui-toolkit/ansi"

type KeyHandler = func(e *KeyEvent)

func OnKey(s *Scope, fn KeyHandler) {
	s.keyHandlers = append(s.keyHandlers, fn)
}

type KeyEvent struct {
	R       rune
	Code    ansi.Key
	stopped bool
}

func (ke *KeyEvent) StopPropagation() {
	ke.stopped = true
}
