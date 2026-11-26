package hook

type EffectRef struct {
	started bool
	cleanup EffectCleanup
	deps    []any
}

func (ref *EffectRef) DepChanged(deps []any) bool {
	if len(ref.deps) != len(deps) {
		return true
	}

	for i := 0; i < len(ref.deps); i += 1 {
		if ref.deps[i] != deps[i] {
			return true
		}
	}

	return false
}

func (ref *EffectRef) Cleanup() {
	if ref.cleanup != nil {
		ref.cleanup()
	}
}

type EffectRefs struct {
	refs   []*EffectRef
	cursor int
}

func (s *EffectRefs) Reset() {
	s.cursor = 0
}

func (s *EffectRefs) NextRef() *EffectRef {
	if len(s.refs) <= s.cursor {
		ref := &EffectRef{}
		s.refs = append(s.refs, ref)
	}

	slot := s.refs[s.cursor]
	s.cursor += 1

	return slot
}

type EffectCleanup = func()
type EffectFunc func() EffectCleanup

func Effect(s *Scope, fn EffectFunc, deps ...any) {
	ref := s.effectRefs.NextRef()

	if ref.started && !ref.DepChanged(deps) {
		return
	}

	if ref.cleanup != nil {
		ref.cleanup()
	}

	ref.cleanup = fn()
	ref.deps = deps
	ref.started = true
}
