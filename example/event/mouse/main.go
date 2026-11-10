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
		//termio.WithMouseSGR,
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
			fmt.Fprintf(term, "%b -> X: %d Y: %d\r\n", token, mouse.X, mouse.Y)
			continue
		}

		fmt.Fprintf(term, "%b\r\n", token)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
