package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/mpgasior/tui-go/bufiox"
	"github.com/mpgasior/tui-go/termx"
)

func main() {
	input, _ := termx.NewTerminalInput(os.Stdin)
	defer input.Close()

	restore, _ := input.MakeRaw()
	defer restore()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := bufiox.NewContextReader(input)
	ctrlC := []byte{0x03}

	var buffer []byte
	const peekSize = 64

	for {
		if ctx.Err() != nil {
			break
		}

		shortCtx, shortCancel := context.WithTimeout(ctx, 20*time.Millisecond)
		b, _ := reader.PeekContext(shortCtx, peekSize)
		shortCancel()
		if len(b) == 0 {
			buffer = nil
			continue
		}

		reader.DiscardContext(ctx, len(b))
		buffer = append(buffer, b...)

		fmt.Printf("[% X]\r\n", buffer)

		if slices.Equal(buffer, ctrlC) {
			cancel()
		}
	}
}
