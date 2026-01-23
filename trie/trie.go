package trie

import (
	"fmt"
	"iter"
)

var ErrAlreadyExists = fmt.Errorf("key already exists in the trie")

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

func (t *Trie[K, V]) Insert(seq iter.Seq[K], value V) error {
	current := t
	for k := range seq {
		next, ok := current.children[k]
		if !ok {
			next = NewTrie[K, V]()
			current.children[k] = next
		}

		current = next
	}

	if current.isTerm {
		return ErrAlreadyExists
	}

	current.value = value
	current.isTerm = true

	return nil
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
