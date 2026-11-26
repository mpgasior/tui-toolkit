package hook

import (
	"sync"
)

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

func (s *StateRefs) NextRef(init func() any) *StateRef {
	if len(s.refs) <= s.cursor {
		val := init()
		s.refs = append(s.refs, &StateRef{val: val})
	}

	slot := s.refs[s.cursor]
	s.cursor += 1

	return slot
}

type SetStateFunc[T any] func(T) T

func SetState[T any](fn func(SetStateFunc[T]), val T) {
	fn(func(T) T {
		return val
	})
}

func State[T any](s *Scope, init T) (state T, setState func(SetStateFunc[T])) {
	return StateFn(s, func() T { return init })
}

func StateFn[T any](s *Scope, init func() T) (state T, setState func(SetStateFunc[T])) {
	ref := s.stateRefs.NextRef(func() any { return init() })

	set := func(fn SetStateFunc[T]) {
		ref.mu.Lock()
		defer ref.mu.Unlock()

		ref.val = fn(ref.val.(T))
		s.schedule()
	}

	return ref.val.(T), set
}
