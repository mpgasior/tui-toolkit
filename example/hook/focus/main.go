package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mpgasior/tui-toolkit/ansi"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/event"
	"github.com/mpgasior/tui-toolkit/gfx"
	"github.com/mpgasior/tui-toolkit/hook"
	"github.com/mpgasior/tui-toolkit/seq"
	"github.com/mpgasior/tui-toolkit/termio"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	counter := func(title string) hook.Component {
		return func(s *hook.Scope, sf *draw.Surface) {
			state, setState := hook.State(s, 42)

			incr := func(delta int) hook.SetStateFn[int] {
				return func(i int) int { return i + delta }
			}

			s.OnKey(func(e *hook.KeyEvent) {
				switch e.Code {
				case ansi.KeyUp:
					setState(incr(1))
					e.StopPropagation()
				case ansi.KeyDown:
					setState(incr(-1))
					e.StopPropagation()
				}
			})

			boxOptions := draw.BoxOptions{
				Border: draw.BoxBorderThin,
				Style:  gfx.StyleDefault,
			}
			if s.IsFocused() {
				boxOptions.Style = gfx.StyleDefault.SetFgColor(gfx.Green)
			}
			draw.Box(sf, boxOptions)
			draw.String(sf.Inset(0, 0, 0, 2), title, draw.WithTextWriterStyle(boxOptions.Style))

			w := draw.NewTextWriter(sf.Inset(1, 1, 1, 1))
			fmt.Fprintf(w, "counter = %d", state)
		}
	}

	grid := func(s *hook.Scope, sf *draw.Surface) {
		w, h := sf.Size()
		lw, rw := w/2, w-w/2
		th, bh := h/2, h-h/2

		topLeft := sf.Sub(0, 0, lw, th)
		topRight := sf.Sub(lw, 0, rw, th)
		bottomLeft := sf.Sub(0, th, lw, bh)
		bottomRight := sf.Sub(lw, th, rw, bh)

		counter1 := s.Mount("top-left", topLeft, counter("[1]"))
		counter2 := s.Mount("top-right", topRight, counter("[2]"))
		counter3 := s.Mount("bottom-left", bottomLeft, counter("[3]"))
		counter4 := s.Mount("bottom-right", bottomRight, counter("[4]"))

		s.OnKey(func(e *hook.KeyEvent) {
			switch e.R {
			case '1':
				counter1.Focus()
			case '2':
				counter2.Focus()
			case '3':
				counter3.Focus()
			case '4':
				counter4.Focus()
			}
		})
	}

	term, err := termio.OpenSession(os.Stdin, os.Stdout, termio.WithAltScreen)
	if err != nil {
		return err
	}
	defer term.Close()

	w, h, err := term.GetSize()
	if err != nil {
		return err
	}

	screen := gfx.NewScreen(w, h, nil)
	escReader := seq.NewEscReader(context.TODO(), term.ReadContext)
	scanner := bufio.NewScanner(escReader)
	scanner.Split(seq.Scan)
	keyDecoder := event.NewKeyDecoder()

	root := hook.MountRoot(grid)

	for {
		sf := draw.NewSurface(screen.Back)
		hook.Render(root, sf)
		screen.Swap()
		if _, err := screen.WriteTo(term); err != nil {
			return err
		}

		if scanner.Scan() {
			token := scanner.Bytes()
			if key, ok := keyDecoder.Decode(token); ok {
				if key.Code == ansi.KeyCtrlC {
					return nil
				}
				event := &hook.KeyEvent{
					Code: key.Code,
					R:    key.R,
				}
				hook.DispatchKey(root, event)
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}
}
