package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/nimelo/tui-go/bufiox"
	"github.com/nimelo/tui-go/termx"
)

func main() {
	tty, _ := termx.OpenTTY()
	defer tty.Close()

	terminal, _ := termx.NewTerminal(tty.In, tty.Out)
	defer terminal.Close()

	restoreInput, _ := terminal.MakeRaw()
	defer restoreInput()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := bufiox.NewContextReader(terminal)
	ctrlC := []byte{0x03}

	var buffer []byte

	for {
		if ctx.Err() != nil {
			break
		}

		shortCtx, shortCancel := context.WithTimeout(ctx, 20*time.Millisecond)
		b, err := reader.ReadByteContext(shortCtx)
		if errors.Is(err, context.DeadlineExceeded) {
			buffer = nil
			continue
		}
		shortCancel()

		buffer = append(buffer, b)

		fmt.Fprintf(terminal, "[% X]\r\n", buffer)

		if slices.Equal(buffer, ctrlC) {
			cancel()
		}
	}
}
