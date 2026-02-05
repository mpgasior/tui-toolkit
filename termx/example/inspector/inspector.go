package main

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/mpgasior/tui-toolkit/termx"
)

func main() {
	tty, _ := termx.OpenTTY()
	defer tty.Close()

	terminal, _ := termx.New(tty.In, tty.Out)
	defer terminal.Close()

	restore, _ := terminal.MakeRaw()
	defer restore()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buffer := make([]byte, 1024)
	ctrlC := []byte{0x03}

	var wg sync.WaitGroup
	wg.Go(func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-terminal.ResizeC():
				w, h, _ := terminal.GetSize()
				fmt.Fprintf(terminal, "New terminal size is: (%d, %d)\r\n", w, h)
			}
		}
	})

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

	wg.Wait()
}
