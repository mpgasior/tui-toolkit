package hook

import (
	"context"
	"slices"
	"sync"
)

type Scope struct {
	parent   *Scope
	children map[string]*Scope
	mounted  []string

	focusMgr *FocusManager

	keyHandler KeyHandler

	providerRefs ProviderRefs
	stateRefs    *StateRefs
	effectRefs   *EffectRefs

	context context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup

	schedule func()
}

func NewScope(parent *Scope) *Scope {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Scope{
		parent:   parent,
		children: make(map[string]*Scope),

		providerRefs: ProviderRefs{},
		stateRefs:    &StateRefs{},
		effectRefs:   &EffectRefs{},

		context: ctx,
		cancel:  cancel,
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
	s.cancel()
	s.wg.Wait()

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
