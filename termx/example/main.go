package main

import (
	"context"
	"fmt"
	"sync"
	"time"

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
		for {
			peekContext, peekCancel := context.WithTimeout(ctx, 10*time.Second)

			ok, _ := terminal.PeekContext(peekContext)
			peekCancel()

			if peekContext.Err() != nil {
				println("peek was cancelled")
			}

			if ctx.Err() != nil {
				break
			}

			if !ok {
				fmt.Fprintf(tty.Out, "Nothing has been pressed...\r\n")
			}

			if ok {
				fmt.Fprintf(tty.Out, "Read to read ...\r\n")
				buffer := make([]byte, 1024)
				n, err := terminal.ReadContext(ctx, buffer)
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
	})

	timeout := 15 * time.Second
	shutdown := time.After(timeout)
	fmt.Fprintf(tty.Out, "Stopping app in %f seconds... (unless 'q' is pressed)\r\n", timeout.Seconds())

MainLoop:
	for {
		select {
		case <-ctx.Done():
			break MainLoop
		case <-time.After(time.Second):
			fmt.Fprintf(tty.Out, "tick..\r\n")
		case <-shutdown:
			fmt.Fprintf(tty.Out, "shutdown..\r\n")
			cancel()
			break MainLoop
		case b := <-ch:
			fmt.Fprintf(tty.Out, "You've clicked: %c (%v)\r\n", rune(b), b)
			if b == '\x1b' {
				cancel()
			}
		}
	}

	wg.Wait()
	fmt.Fprintf(tty.Out, "closing the app in 1s ..\r\n")
	<-time.After(time.Second)
}
