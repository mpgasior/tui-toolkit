package hook

import (
	"maps"

	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/event"
)

type KeyHandler = func(e *KeyEvent)
type MouseHandler = func(e event.Mouse)

type Scope struct {
	component Component
	parent    *Scope

	focus *Focus

	keyHandler   KeyHandler
	mouseHandler MouseHandler

	stateRefs *StateRefs

	children map[string]*Scope
}

func NewScope(c Component) *Scope {
	s := &Scope{
		component: c,
		children:  make(map[string]*Scope),
		stateRefs: &StateRefs{},
	}

	return s
}

func (s *Scope) Render(sf *draw.Surface) {
	s.stateRefs.Reset()
	s.component(s, sf)
}

func (s *Scope) OnKey(fn KeyHandler) {
	s.keyHandler = fn
}

func (s *Scope) OnMouse(fn MouseHandler) {
	s.mouseHandler = fn
}

func (s *Scope) IsFocusable() bool {
	return s.keyHandler != nil || s.mouseHandler != nil
}

func (s *Scope) Focus() {
	s.focus.SetFocus(s)
}

func (s *Scope) IsFocused() bool {
	return s.focus.IsFocused(s)
}

func (s *Scope) IsInFocusChain() bool {
	return s.focus.IsInChain(s)
}

func (s *Scope) Mount(id string, sf *draw.Surface, c Component) *Scope {
	if scope, ok := s.children[id]; ok {
		prev := maps.Clone(scope.children)
		scope.Render(sf)
		for prevId, prevScope := range prev {
			if _, ok := scope.children[prevId]; !ok {
				prevScope.unmount()
			}
		}

		return scope
	}

	scope := NewScope(c)
	scope.parent = s
	scope.focus = s.focus
	s.children[id] = scope
	scope.Render(sf)

	return scope
}

func (s *Scope) unmount() {
	s.keyHandler = nil
	s.mouseHandler = nil
	s.stateRefs = nil

	for _, scope := range s.children {
		scope.unmount()
	}
	s.children = nil

	if s.focus.IsFocused(s) {
		s.focus.SetFocus(s.parent)
	}

	s.parent = nil
}
