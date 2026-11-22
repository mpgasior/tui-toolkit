package hook

import (
	"github.com/mpgasior/tui-toolkit/draw"
)

type Component func(s *Scope, sf *draw.Surface)

type Root struct {
	*Scope
}

func Render(s *Root, sf *draw.Surface) {
	s.Render(sf)
}

func DispatchKey(s *Root, e *KeyEvent) {
	for cursor := s.focus.focused; cursor != nil; cursor = cursor.parent {
		if cursor.keyHandler != nil {
			cursor.keyHandler(e)

			if e.stopped {
				return
			}
		}
	}
}

func MountRoot(c Component) *Root {
	s := NewScope(c)
	s.focus = &Focus{focused: s}
	r := &Root{
		Scope: s,
	}

	return r
}
