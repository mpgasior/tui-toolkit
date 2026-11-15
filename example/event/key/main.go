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
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	term, err := termio.OpenSession(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}
	defer term.Close()

	keyDecoder := event.NewKeyDecoder()

	ctx := context.Background()
	r := seq.NewEscReader(ctx, term.ReadContext)
	scanner := bufio.NewScanner(r)
	scanner.Split(seq.Scan)

	for scanner.Scan() {
		token := scanner.Bytes()

		if key, ok := keyDecoder.Decode(token); ok {
			fmt.Fprintf(term, "%b -> Code: %v Rune: %c\r\n", token, key.Code, key.R)
			if key.Code == ansi.KeyCtrlC {
				return nil
			}
			continue
		}

		fmt.Fprintf(term, "%b\r\n", token)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
