package scope

import "errors"

type Step[T any] struct {
	Apply   func(T) error
	Restore func(T) error
}

type Scope[T any] struct {
	val          T
	steps        []Step[T]
	appliedCount int
}

func New[T any](val T, steps ...Step[T]) *Scope[T] {
	return &Scope[T]{
		val:   val,
		steps: steps,
	}
}

func (s *Scope[T]) Apply() error {
	for _, step := range s.steps {
		if err := step.Apply(s.val); err != nil {
			restoreErr := s.Restore()
			return errors.Join(err, restoreErr)
		}
		s.appliedCount += 1
	}
	return nil
}

func (s *Scope[T]) Restore() error {
	var errs []error
	for i := s.appliedCount - 1; i >= 0; i -= 1 {
		if err := s.steps[i].Restore(s.val); err != nil {
			errs = append(errs, err)
		}
	}

	s.appliedCount = 0
	return errors.Join(errs...)
}
