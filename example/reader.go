package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mpgasior/tui-toolkit/termio"
	"golang.org/x/term"
)

func main() {
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(fd, state)

	r, err := termio.NewReader(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctrlC := [8]byte{3}
	for {
		if ctx.Err() != nil {
			break
		}

		var buf [8]byte
		readCtx, stop := context.WithTimeout(ctx, time.Second*3)
		n, err := r.ReadContext(readCtx, buf[:])
		stop()

		if err != nil {
			fmt.Printf("%v\r\n", err)
			continue
		}

		if n == 0 {
			continue
		}

		fmt.Printf("%v\r\n", buf)

		if buf == ctrlC {
			break
		}
	}
}
