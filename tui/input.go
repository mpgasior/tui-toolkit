package tui

import (
	"context"
	"slices"
	"unicode/utf8"

	"github.com/mpgasior/tui-go/termx"
	"github.com/mpgasior/tui-go/trie"
	"github.com/mpgasior/tui-go/vt"
)

func Input(terminal *termx.Terminal) func(ctx context.Context, ch chan<- Event) {
	return func(ctx context.Context, ch chan<- Event) {
		trie := trie.New[byte, vt.Key]()
		for seq, key := range vt.SequenceToKey {
			_ = trie.Insert(slices.Values([]byte(seq)), key)
		}

		scanner := vt.NewSequenceScanner(terminal, vt.ScanInitial)
		for scanner.ScanContext(ctx) {
			seq := scanner.Sequence()
			if seq.Is(vt.SeqPaste) {
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

			if key, ok := trie.Get(slices.Values(seq.Data)); ok {
				e := KeyEvent{Key: key}
				select {
				case <-ctx.Done():
					return
				case ch <- e:
				}
				continue
			}

			if seq.Is(vt.SeqUTF8) {
				r, _ := utf8.DecodeRune(seq.Data)
				e := KeyEvent{Rune: r}
				select {
				case <-ctx.Done():
					return
				case ch <- e:
				}
			}
		}
	}
}
