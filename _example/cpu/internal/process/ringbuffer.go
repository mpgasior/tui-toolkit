package process

type RingBuffer[T any] struct {
	data []T
	head int
	full bool
}

func NewRingBuffer[T any](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data: make([]T, size),
	}
}

func (r *RingBuffer[T]) Full() bool {
	return r.full
}

func (r *RingBuffer[T]) Len() int {
	if r.full {
		return len(r.data)
	}

	return r.head
}

func (r *RingBuffer[T]) Get(i int) (T, bool) {
	if i < 0 || i >= r.Len() {
		var zero T
		return zero, false
	}

	if !r.full {
		return r.data[i], true
	}

	idx := (r.head + i) % len(r.data)
	return r.data[idx], true
}

func (r *RingBuffer[T]) Push(val T) {
	r.data[r.head] = val
	r.head = (r.head + 1) % len(r.data)

	if r.head == 0 {
		r.full = true
	}
}

func (r *RingBuffer[T]) All() []T {
	if !r.full {
		out := make([]T, r.head)
		copy(out, r.data)
		return out
	}

	out := make([]T, len(r.data))
	n := copy(out, r.data[r.head:])
	copy(out[n:], r.data[:r.head])

	return out
}
