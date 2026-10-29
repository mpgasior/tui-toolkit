package scope

type Builder[T any] struct {
	steps []Step[T]
	val   T
}

func NewBuilder[T any](val T, steps ...Step[T]) *Builder[T] {
	return &Builder[T]{val: val, steps: steps}
}

func (b *Builder[T]) AddStep(s Step[T]) *Builder[T] {
	b.steps = append(b.steps, s)
	return b
}

func (b *Builder[T]) Add(apply func(T) error, restore func(T) error) *Builder[T] {
	step := Step[T]{
		Apply:   apply,
		Restore: restore,
	}

	return b.AddStep(step)
}

func (b *Builder[T]) Scope() *Scope[T] {
	return New(b.val, b.steps...)
}
