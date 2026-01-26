package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"unicode"
	"unicode/utf8"

	"github.com/mpgasior/tui-go/termx"
	"github.com/mpgasior/tui-go/vt"
)

func main() {
	terminal, _ := termx.New(os.Stdin, os.Stdout)
	defer terminal.Close()

	restoreInput, _ := terminal.MakeRaw()
	defer restoreInput()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exitMode, _ := vt.EnterMode(terminal, vt.ModeBracketedPaste)
	defer exitMode()
	io.WriteString(terminal, vt.QueryTerminalName)
	io.WriteString(terminal, vt.QueryBgColor)
	io.WriteString(terminal, vt.QueryFgColor)
	io.WriteString(terminal, vt.QueryCursorColor)

	scanner := vt.NewSequenceScanner(terminal, vt.ScanInitial)
	ctrlC := []byte{0x03}

	for scanner.ScanContext(ctx) {
		seq := scanner.Sequence()

		fmt.Fprintf(terminal, "%s: [% X]", seq.Type.String(), seq.Data)

		if seq.Is(vt.SeqUTF8) {
			r, _ := utf8.DecodeRune(seq.Data)
			if unicode.IsPrint(r) {
				fmt.Fprintf(terminal, " (%c)", r)
			}
		}

		if seq.Is(vt.SeqPaste) {
			str := string(seq.Data)
			if utf8.ValidString(str) {
				fmt.Fprintf(terminal, " %s", str)
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
