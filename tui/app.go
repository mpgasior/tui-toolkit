package tui

import (
	"context"
	"slices"
	"unicode/utf8"

	"github.com/mpgasior/tui-go/termx"
	"github.com/mpgasior/tui-go/trie"
	"github.com/mpgasior/tui-go/vt"
)

type App struct {
}

func (a App) Run(c Component) error {
	tty, err := termx.OpenTTY()
	if err != nil {
		return err
	}
	defer tty.Close()

	terminal, err := termx.New(tty.In, tty.Out)
	if err != nil {
		return err
	}
	defer terminal.Close()

	restore, err := terminal.MakeRaw()
	if err != nil {
		return err
	}
	defer restore()

	w, h, err := terminal.GetSize()
	if err != nil {
		return err
	}

	ch := make(chan Event)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runtime := NewRuntime(ch)
	dispatch, shutdown := runtime.Start(ctx)
	defer shutdown()

	dispatch(c.Init())

	reader := func(ctx context.Context, ch chan<- Event) {
		trie := trie.New[byte, vt.Key]()
		for seq, key := range vt.SequenceToKey {
			_ = trie.Insert(slices.Values([]byte(seq)), key)
		}

		scanner := vt.NewSequenceScanner(terminal, vt.ScanInitial)
		for scanner.ScanContext(ctx) {
			seq := scanner.Sequence()
			if seq.Is(vt.SeqPaste) {
				var bytes []byte
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

	dispatch(TaskF(reader))

	renderer := NewRenderer(w, h)

	for {
		c.Render(RenderContext{
			Viewport: NewViewport(renderer.Back),
			Focused:  true,
		})

		renderer.SwapBuffers()
		renderer.WriteTo(terminal)

		select {
		case <-ctx.Done():
			return nil
		case ev := <-ch:
			if ev == ShutdownEvent {
				cancel()
				return nil
			}

			task := c.Update(ev)
			dispatch(task)
		}
	}
}
