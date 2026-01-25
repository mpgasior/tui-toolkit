package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"

	"github.com/mpgasior/tui-go/termx"
	"github.com/mpgasior/tui-go/trie"
	"github.com/mpgasior/tui-go/vt"
)

func main() {
	tty, _ := termx.OpenTTY()
	terminal, _ := termx.NewTerminal(tty.In, tty.Out)
	defer terminal.Close()

	restoreInput, _ := terminal.MakeRaw()
	defer restoreInput()

	exitAltScreen, _ := vt.EnterMode(os.Stdout, vt.ModeAlternateScreen)
	defer exitAltScreen()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scanner := vt.NewSequenceScanner(terminal, vt.ScanInitial)

	trie := trie.NewTrie[byte, vt.Key]()
	for seq, key := range vt.SequenceToKey {
		_ = trie.Insert(slices.Values([]byte(seq)), key)
	}

	for scanner.ScanContext(ctx) {
		seq := scanner.Sequence()

		fmt.Fprintf(terminal, "%s: [% X]\r\n", seq.Type.String(), seq.Data)

		if key, ok := trie.Get(slices.Values(seq.Data)); ok {
			if key == vt.KeyCtrlC {
				cancel()
			}

			if key == vt.KeyCtrlE {
				exitAltScreen()
				restoreInput()

				cmd := exec.Command("nvim")
				cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, tty.In, tty.Out
				cmd.Env = os.Environ()

				if err := cmd.Run(); err != nil {
					fmt.Fprintf(terminal, "Could not start nvim due to: %v\r\n", err)
				}

				exitAltScreen, _ := vt.EnterMode(terminal, vt.ModeAlternateScreen)
				defer exitAltScreen()

				restoreInput, _ := terminal.MakeRaw()
				defer restoreInput()

				fmt.Fprintf(terminal, "Finished with %v\r\n", cmd.ProcessState)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %v\r\n", err)
	}
}
