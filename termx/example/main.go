package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"unicode"

	"github.com/nimelo/tui-go/termx"
)

func main() {
	tty, err := termx.OpenTTY()
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	terminal, err := termx.NewTerminal(tty.In, tty.Out)
	if err != nil {
		panic(err)
	}
	defer terminal.Close()

	restoreOutput, err := terminal.EnterAltScreen()
	if err != nil {
		panic(err)
	}
	defer restoreOutput()

	restoreInput, err := terminal.MakeRaw()
	if err != nil {
		panic(err)
	}
	defer restoreInput()

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan byte)
	var wg sync.WaitGroup

	wg.Go(func() {
		read(ctx, terminal, ch)
	})

	timeout := 15 * time.Second
	shutdown := time.After(timeout)
	fmt.Fprintf(terminal, "Stopping app in %f seconds... (unless 'ESC' is pressed)\r\n", timeout.Seconds())
	terminal.Flush()

MainLoop:
	for {
		select {
		case <-ctx.Done():
			break MainLoop
		case <-time.After(time.Second):
			fmt.Fprintf(terminal, "tick..\r\n")
			terminal.Flush()
		case <-shutdown:
			fmt.Fprintf(terminal, "shutdown..\r\n")
			terminal.Flush()
			cancel()
			break MainLoop
		case b := <-ch:
			r := rune(b)
			if unicode.IsPrint(r) {
				fmt.Fprintf(terminal, "You've clicked: %c (%v)\r\n", rune(b), b)
			} else {
				fmt.Fprintf(terminal, "You've clicked: 0x%0x \r\n", b)
			}

			terminal.Flush()
			if b == '\x1b' {
				cancel()
			}
		}
		terminal.Flush()
	}

	wg.Wait()
	fmt.Fprintf(terminal, "closing the app in 1s ..\r\n")
	terminal.Flush()
	<-time.After(time.Second)
}

func read(ctx context.Context, terminal termx.Terminal, ch chan<- byte) {
	for {
		readyContext, readyCancel := context.WithTimeout(ctx, 3*time.Second)

		readyErr := terminal.Ready(readyContext)
		readyCancel()

		if errors.Is(readyErr, context.DeadlineExceeded) || errors.Is(readyErr, context.Canceled) {
			fmt.Fprint(terminal, "Nothing has been pressed...\r\n")
			terminal.Flush()
		}

		if ctx.Err() != nil {
			fmt.Fprint(terminal, "Breaking main loop... \r\n")
			terminal.Flush()
			break
		}

		if readyErr == nil {
			fmt.Fprint(terminal, "Performing read... \r\n")
			terminal.Flush()

			buffer := make([]byte, 1024)
			n, err := terminal.Read(ctx, buffer)
			if err != nil {
				break
			}

			if n == 0 {
				continue
			}

			for i := range n {
				select {
				case <-ctx.Done():
					return
				case ch <- buffer[i]:
				}
			}
		}
	}
}
