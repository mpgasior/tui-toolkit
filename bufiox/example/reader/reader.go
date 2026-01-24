package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/nimelo/tui-go/bufiox"
	"github.com/nimelo/tui-go/termx"
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

	for {
		if ctx.Err() != nil {
			break
		}

		shortCtx, shortCancel := context.WithTimeout(ctx, 20*time.Millisecond)
		b, err := reader.PeekContext(shortCtx, 1)
		shortCancel()
		if errors.Is(err, context.DeadlineExceeded) {
			buffer = nil
			continue
		}

		_, _ = reader.Discard(1)
		buffer = append(buffer, b...)

		fmt.Printf("[% X]\r\n", buffer)

		if slices.Equal(buffer, ctrlC) {
			cancel()
		}
	}
}
