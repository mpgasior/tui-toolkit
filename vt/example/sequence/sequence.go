package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"unicode"
	"unicode/utf8"

	"github.com/nimelo/tui-go/termx"
	"github.com/nimelo/tui-go/vt"
)

func main() {
	terminal, _ := termx.NewTerminal(os.Stdin, os.Stdout)
	defer terminal.Close()

	restoreInput, _ := terminal.MakeRaw()
	defer restoreInput()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	terminal.Write([]byte(vt.QueryTerminalName))
	terminal.Write([]byte(vt.QueryBgColor))
	terminal.Write([]byte(vt.QueryFgColor))
	terminal.Write([]byte(vt.QueryCursorColor))

	scanner := vt.NewSequenceScanner(terminal, vt.ScanInitial)
	ctrlC := []byte{0x03}

	for scanner.ScanContext(ctx) {
		seq := scanner.Sequence()

		fmt.Fprintf(terminal, "%s: [% X]", seq.Type.String(), seq.Data)

		if seq.Is(vt.SeqUtf8) {
			r, _ := utf8.DecodeRune(seq.Data)
			if unicode.IsPrint(r) {
				fmt.Fprintf(terminal, " (%c)", r)
			}
		}

		if slices.Equal(seq.Data, ctrlC) {
			cancel()
		}

		fmt.Fprintf(terminal, "\r\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(terminal, "Error: %v\r\n", err)
	}
}
