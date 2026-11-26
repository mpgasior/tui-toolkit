package hook

import (
	"slices"
)

type Scope struct {
	parent   *Scope
	children map[string]*Scope
	mounted  []string

	focusMgr *FocusManager

	keyHandler KeyHandler
	stateRefs  *StateRefs
	effectRefs *EffectRefs

	schedule func()
}

func NewScope(parent *Scope) *Scope {
	s := &Scope{
		parent:   parent,
		children: make(map[string]*Scope),

		stateRefs:  &StateRefs{},
		effectRefs: &EffectRefs{},
	}

	if parent != nil {
		s.focusMgr = parent.focusMgr
		s.schedule = parent.schedule
	}

	return s
}

func (s *Scope) Render(comp Component) RenderFunc {
	s.renderBegin()
	render := comp(s)
	s.renderEnd()

	return render
}

func (s *Scope) renderBegin() {
	s.stateRefs.Reset()
	s.effectRefs.Reset()
	s.mounted = s.mounted[:0]
}

func (s *Scope) renderEnd() {
	for k, v := range s.children {
		if !slices.Contains(s.mounted, k) {
			v.Unmount()
			delete(s.children, k)
		}
	}
}

func (s *Scope) Mount(id string, comp Component) RenderFunc {
	target, ok := s.children[id]
	if !ok {
		target = NewScope(s)
		s.children[id] = target
	}

	s.mounted = append(s.mounted, id)

	return target.Render(comp)
}

func (s *Scope) Unmount() {
	s.keyHandler = nil
	s.stateRefs = nil

	for _, scope := range s.children {
		scope.Unmount()
	}
	s.children = nil

	if s.focusMgr.IsFocused(s) {
		s.focusMgr.SetFocus(s.parent)
	}

	s.parent = nil
}
