package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mpgasior/tui-toolkit/ansi"
	"github.com/mpgasior/tui-toolkit/event"
	"github.com/mpgasior/tui-toolkit/seq"
	"github.com/mpgasior/tui-toolkit/termio"
)

func main() {
	mods := []termio.SessionModifier{
		termio.WithAltScreen,
		termio.WithMouseNormal,
		termio.WithMouseAny,
	}
	term, err := termio.OpenSession(os.Stdin, os.Stdout, mods...)
	if err != nil {
		log.Fatalln(err)
	}
	defer term.Close()

	keyDecoder := event.NewKeyDecoder()
	mouseDecoder := event.NewMouseDecoder()
	ctx := context.Background()
	r := seq.NewEscReader(ctx, term.ReadContext)
	scanner := bufio.NewScanner(r)
	scanner.Split(seq.Scan)

	for scanner.Scan() {
		token := scanner.Bytes()

		if key, ok := keyDecoder.Decode(token); ok {
			fmt.Fprintf(term, "%b -> Code: %v Rune: %c\r\n", token, key.Code, key.R)
			if key.Code == ansi.KeyCtrlC {
				break
			}
			continue
		}

		if mouse, ok := mouseDecoder.Decode(token); ok {
			var button, action string
			var mod []string
			switch mouse.Button {
			case event.MouseButtonLeft:
				button = "left"
			case event.MouseButtonMiddle:
				button = "middle"
			case event.MouseButtonRight:
				button = "right"
			case event.MouseButtonWheelUp:
				button = "wheel up"
			case event.MouseButtonWheelDown:
				button = "wheel down"
			default:
				button = "unknown"
			}

			switch mouse.Action {
			case event.MouseActionMove:
				action = "move"
			case event.MouseActionPress:
				action = "press"
			case event.MouseActionRelease:
				action = "release"
			}

			if mouse.Mods&event.MouseModifierCtrl != 0 {
				mod = append(mod, "ctrl")
			}
			if mouse.Mods&event.MouseModifierAlt != 0 {
				mod = append(mod, "alt")
			}
			if mouse.Mods&event.MouseModifierShift != 0 {
				mod = append(mod, "shift")
			}

			fmt.Fprintf(term, "%b -> %s %s %v X: %d Y: %d\r\n", token, button, action, mod, mouse.X, mouse.Y)
			continue
		}

		fmt.Fprintf(term, "%b\r\n", token)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
