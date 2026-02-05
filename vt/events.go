package vt

import (
	"context"
	"slices"
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/iox"
	"github.com/mpgasior/tui-toolkit/trie"
)

type KeyEvent struct {
	Key  Key
	Rune rune
}

func (e KeyEvent) IsKey(keys ...Key) bool {
	return slices.Contains(keys, e.Key)
}

func (e KeyEvent) IsRune(r rune) bool {
	return e.Rune == r
}

type PasteEvent struct {
	Bytes []byte
}

func Events(reader iox.ContextReader) (func(ctx context.Context), <-chan any) {
	ch := make(chan any)

	f := func(ctx context.Context) {
		trie := trie.New[byte, Key]()
		for seq, key := range SequenceToKey {
			_ = trie.Insert(slices.Values([]byte(seq)), key)
		}

		decodeKey := func(seq Sequence) (key Key) {
			key, _ = trie.Get(slices.Values(seq.Data))
			return key
		}

		decodeRune := func(seq Sequence) (r rune) {
			r = utf8.RuneError
			if ok := seq.Is(SeqUTF8); ok {
				r, _ = utf8.DecodeRune(seq.Data)
			}
			return r
		}

		scanner := NewSequenceScanner(reader, ScanInitial)
		for scanner.ScanContext(ctx) {
			seq := scanner.Sequence()
			if seq.Is(SeqPaste) {
				bytes := make([]byte, len(seq.Data))
				copy(bytes, seq.Data)
				e := PasteEvent{Bytes: bytes}
				select {
				case <-ctx.Done():
					return
				case ch <- e:
				}
				continue
			}

			key := decodeKey(seq)
			r := decodeRune(seq)
			e := KeyEvent{Key: key, Rune: r}
			select {
			case <-ctx.Done():
				return
			case ch <- e:
			}
		}
	}

	return f, ch
}
