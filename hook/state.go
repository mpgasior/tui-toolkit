package hook

import "sync"

type StateRef struct {
	mu  sync.Mutex
	val any
}

type StateRefs struct {
	refs   []*StateRef
	cursor int
}

func (s *StateRefs) Reset() {
	s.cursor = 0
}

func (s *StateRefs) NextRef(init any) *StateRef {
	if len(s.refs) <= s.cursor {
		s.refs = append(s.refs, &StateRef{val: init})
	}

	slot := s.refs[s.cursor]
	s.cursor += 1

	return slot
}

type SetStateFn[T any] func(T) T

func State[T any](s *Scope, init T) (state T, setState func(SetStateFn[T])) {
	ref := s.stateRefs.NextRef(init)

	set := func(fn SetStateFn[T]) {
		ref.mu.Lock()
		defer ref.mu.Unlock()

		ref.val = fn(ref.val.(T))
	}

	return ref.val.(T), set
}
