package event

import (
	"slices"
	"unicode/utf8"

	"github.com/mpgasior/tui-toolkit/ansi"
	"github.com/mpgasior/tui-toolkit/internal/trie"
)

type Key struct {
	Code ansi.Key
	R    rune
}

type KeyDecoder struct {
	keys *trie.Trie[byte, Key]
}

func NewKeyDecoder() *KeyDecoder {
	d := &KeyDecoder{
		keys: trie.New[byte, Key](),
	}

	for k, v := range ansi.SequenceToKey {
		buf := []byte(k)

		key := Key{Code: v}
		if r, ok := d.decodeRune(buf); ok {
			key.R = r
		}

		d.RegisterBinding(buf, key)
	}

	return d
}

func (d *KeyDecoder) RegisterBinding(seq []byte, k Key) {
	d.keys.Insert(slices.Values(seq), k)
}

func (d *KeyDecoder) decodeRune(seq []byte) (r rune, ok bool) {
	if r, size := utf8.DecodeRune(seq); r != utf8.RuneError && size == len(seq) {
		return r, true
	}
	return 0, false
}

func (d *KeyDecoder) Decode(seq []byte) (Key, bool) {
	if key, ok := d.keys.Get(slices.Values(seq)); ok {
		return key, ok
	}

	if r, ok := d.decodeRune(seq); ok {
		return Key{Code: ansi.KeyUnknown, R: r}, true
	}

	var zero Key
	return zero, false
}
