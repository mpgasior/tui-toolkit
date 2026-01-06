package main

import (
	"context"
	"fmt"
	"io"
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
		a := make([]byte, 1)
		for {
			n, err := terminal.ReadContext(ctx, a)

			if err != nil {
				break
			}

			if n == 0 {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case ch <- a[0]:
			}
		}
	})

	shutdown := time.After(5 * time.Second)
	io.WriteString(tty.Out, "Stopping app in 5 seconds... (unless 'q' is pressed)\r\n")

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
			if b == 'q' {
				cancel()
			}
		}
	}

	wg.Wait()
	fmt.Fprintf(tty.Out, "closing the app in 1s ..\r\n")
	<-time.After(time.Second)
}
