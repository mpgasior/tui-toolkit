package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"unicode"
	"unicode/utf8"

	"github.com/nimelo/tui-go/termx"
	"github.com/nimelo/tui-go/trie"
	"github.com/nimelo/tui-go/vt"
)

func main() {
	terminal, _ := termx.NewTerminal(os.Stdin, os.Stdout)
	defer terminal.Close()

	restoreInput, _ := terminal.MakeRaw()
	defer restoreInput()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scanner := vt.NewSequenceScanner(terminal, vt.ScanInitial)

	trie := trie.NewTrie[byte, vt.Key]()
	for seq, key := range vt.SequenceToKey {
		trie.Insert(slices.Values([]byte(seq)), key)
	}

	for scanner.ScanContext(ctx) {
		seq := scanner.Sequence()

		fmt.Fprintf(terminal, "%s: [% X]", seq.Type.String(), seq.Data)

		if seq.Is(vt.SeqUtf8) {
			r, _ := utf8.DecodeRune(seq.Data)
			if unicode.IsPrint(r) {
				fmt.Fprintf(terminal, " %c", r)
			}
		}

		if key, ok := trie.Get(slices.Values(seq.Data)); ok {
			fmt.Fprintf(terminal, " => % s", key.Equivalents())

			if key == vt.KeyCtrlC {
				cancel()
			}
		}

		fmt.Fprintf(terminal, "\r\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(terminal, "Error: %v\r\n", err)
	}
}
