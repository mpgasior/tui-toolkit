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
	terminal, _ := termx.NewTerminalInput(os.Stdin)
	defer terminal.Close()

	restore, _ := terminal.MakeRaw()
	defer restore()

	exit, _ := vt.EnterMode(os.Stdout, vt.ModeBracketedPaste)
	defer exit()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scanner := vt.NewSequenceScanner(terminal, vt.ScanInitial)

	trie := trie.NewTrie[byte, vt.Key]()
	for seq, key := range vt.SequenceToKey {
		trie.Insert(slices.Values([]byte(seq)), key)
	}

	for scanner.ScanContext(ctx) {
		seq := scanner.Sequence()

		fmt.Printf("%s: [% X]", seq.Type.String(), seq.Data)

		if seq.Is(vt.SeqUTF8) {
			r, _ := utf8.DecodeRune(seq.Data)
			if unicode.IsPrint(r) {
				fmt.Printf(" %c", r)
			}
		}

		if key, ok := trie.Get(slices.Values(seq.Data)); ok {
			fmt.Printf(" => % s", key.Equivalents())

			if key == vt.KeyCtrlC {
				cancel()
			}
		}

		fmt.Printf("\r\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %v\r\n", err)
	}
}
