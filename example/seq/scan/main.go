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
	mods := []termio.SessionModifier{
		termio.WithBracketedPaste,
	}
	term, err := termio.OpenSession(os.Stdin, os.Stdout, mods...)
	if err != nil {
		log.Fatal(err)
	}
	defer term.Close()

	r := seq.NewEscReader(context.Background(), term.ReadContext, 0)
	scanner := bufio.NewScanner(r)
	scanner.Split(seq.Scan)

	for scanner.Scan() {
		token := scanner.Bytes()
		fmt.Fprintf(term, "%v\r\n", token)
		if bytes.Equal(token, []byte{ansi.ByteEsc}) {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(term, "%v\r\n", err)
		if me, ok := errors.AsType[*seq.MalformedSeqErr](err); ok {
			fmt.Fprintf(term, "%v\r\n", me.Data())
		}
	}
}
