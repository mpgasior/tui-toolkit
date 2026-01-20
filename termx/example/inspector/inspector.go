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

	restoreOutput, _ := terminal.EnterAltScreen()
	defer restoreOutput()

	restoreInput, _ := terminal.MakeRaw()
	defer restoreInput()

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

		fmt.Fprintf(terminal, "[% X]\r\n", sequence)

		if slices.Equal(sequence, ctrlC) {
			cancel()
		}
	}
}
