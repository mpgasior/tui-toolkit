package hook

type ProviderRefs map[any]any

type Key[T any] struct {
	name string
}

func NewKey[T any](name string) Key[T] {
	return Key[T]{name}
}

func Provide[T any](s *Scope, key Key[T], val T) {
	s.providerRefs[key] = val
}

func Consume[T any](s *Scope, key Key[T]) (val T, ok bool) {
	for cursor := s; cursor != nil; cursor = cursor.parent {
		if entry, found := cursor.providerRefs[key]; found {
			return entry.(T), found
		}
	}

	return val, ok
}
