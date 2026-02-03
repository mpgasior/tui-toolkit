package tui

import (
	"context"
	"slices"
	"unicode/utf8"

	"github.com/mpgasior/tui-go/termx"
	"github.com/mpgasior/tui-go/trie"
	"github.com/mpgasior/tui-go/vt"
)

func CaptureResize(terminal *termx.Terminal) func(ctx context.Context, ch chan<- Event) {
	return func(ctx context.Context, ch chan<- Event) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-terminal.ResizeC():
				w, h, err := terminal.GetSize()
				if err != nil {
					continue
				}
				ev := ResizeEvent{
					Width: w, Height: h,
				}
				select {
				case <-ctx.Done():
				case ch <- ev:
				}
			}
		}
	}
}

func CaptureInput(terminal *termx.Terminal) func(ctx context.Context, ch chan<- Event) {
	return func(ctx context.Context, ch chan<- Event) {
		trie := trie.New[byte, vt.Key]()
		for seq, key := range vt.SequenceToKey {
			_ = trie.Insert(slices.Values([]byte(seq)), key)
		}

		decodeKey := func(seq vt.Sequence) (key vt.Key) {
			key, _ = trie.Get(slices.Values(seq.Data))
			return key
		}

		decodeRune := func(seq vt.Sequence) (r rune) {
			r = utf8.RuneError
			if ok := seq.Is(vt.SeqUTF8); ok {
				r, _ = utf8.DecodeRune(seq.Data)
			}
			return r
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
}
