package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mpgasior/tui-toolkit/ansi"
	"github.com/mpgasior/tui-toolkit/seq"
	"github.com/mpgasior/tui-toolkit/termio"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	mods := []termio.SessionModifier{
		termio.WithBracketedPaste,
	}
	term, err := termio.OpenSession(os.Stdin, os.Stdout, mods...)
	if err != nil {
		return err
	}
	defer term.Close()

	r := seq.NewEscReader(context.Background(), term.ReadContext)
	scanner := bufio.NewScanner(r)
	scanner.Split(seq.Scan)

	for scanner.Scan() {
		token := scanner.Bytes()
		fmt.Fprintf(term, "%v\r\n", token)
		if bytes.Equal(token, []byte{ansi.EscByte}) {
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(term, "%v\r\n", err)
		if me, ok := errors.AsType[*seq.MalformedSeqErr](err); ok {
			fmt.Fprintf(term, "%v\r\n", me.Data())
		}
	}

	return nil
}
