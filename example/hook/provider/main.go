package main

import (
	"bufio"
	"context"
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

	root := hook.Root(func() {})
	defer root.Unmount()

	for {
		sf := draw.NewSurface(screen.Back)
		render := root.Render(app())
		render(sf)

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

type theme struct {
	Fg gfx.Color
	Bg gfx.Color
}

var themeKey hook.Key[theme]

func app() hook.Component {
	return func(s *hook.Scope) hook.RenderFunc {
		tm, setTheme := hook.State(s, theme{Fg: gfx.Red, Bg: gfx.Green})
		hook.Provide(s, themeKey, tm)

		hook.OnKey(s, func(key *hook.KeyEvent) {
			switch key.Code {
			case ansi.KeyT:
				setTheme(func(prev theme) theme {
					return theme{
						Fg: prev.Bg,
						Bg: prev.Fg,
					}
				})
			}
		})

		tile := s.Mount("tile", tile())
		return func(sf *draw.Surface) {
			tile(sf)
		}
	}
}

func tile() hook.Component {
	return func(s *hook.Scope) hook.RenderFunc {
		theme, ok := hook.Consume(s, themeKey)

		return func(sf *draw.Surface) {
			if !ok {
				draw.String(sf, "theme not found!")
				return
			}

			style := gfx.StyleDefault.
				SetFgColor(theme.Fg).
				SetBgColor(theme.Bg)
			draw.Fill(sf, ' ', style)
		}
	}
}
