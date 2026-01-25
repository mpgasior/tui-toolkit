package main

import (
	"context"
	"fmt"
	"slices"

	"github.com/nimelo/tui-go/termx"
)

func main() {
	tty, _ := termx.OpenTTY()
	defer tty.Close()

	terminal, _ := termx.NewTerminal(tty.In, tty.Out)
	defer terminal.Close()

	restore, _ := terminal.MakeRaw()
	defer restore()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buffer := make([]byte, 1024)
	ctrlC := []byte{0x03}

	for {
		if ctx.Err() != nil {
			break
		}

		n, _ := terminal.ReadContext(ctx, buffer)
		sequence := buffer[:n]

		if n == 0 {
			continue
		}

		fmt.Fprintf(terminal, "[% X]\r\n", sequence)

		if slices.Equal(sequence, ctrlC) {
			cancel()
		}
	}
}
