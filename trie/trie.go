package trie

import "iter"

type Trie[K comparable, V any] struct {
	children map[K]*Trie[K, V]
	value    V
	isTerm   bool
}

func NewTrie[K comparable, V any]() *Trie[K, V] {
	return &Trie[K, V]{
		children: make(map[K]*Trie[K, V]),
	}
}

func (t *Trie[K, V]) Insert(seq iter.Seq[K], value V) {
	current := t
	for k := range seq {
		if _, ok := current.children[k]; !ok {
			current.children[k] = NewTrie[K, V]()
		}

		current = current.children[k]
	}

	current.value = value
	current.isTerm = true
}

func (t *Trie[K, V]) Get(seq iter.Seq[K]) (val V, found bool) {
	current := t

	for k := range seq {
		next, ok := current.children[k]
		if !ok {
			return val, false
		}
		current = next
	}

	if !current.isTerm {
		return val, false
	}

	return current.value, true
}
